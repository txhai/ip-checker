package ipchecker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type CIDRList struct {
	sync.Mutex
	v4 []*net.IPNet
	v6 []*net.IPNet
}

func createCIDRList(fileName string) (*CIDRList, error) {
	l := new(CIDRList)
	l.Lock()
	defer l.Unlock()
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	v4 := make([]*net.IPNet, 0)
	v6 := make([]*net.IPNet, 0)
	scanner := bufio.NewScanner(file)
	fmt.Printf("loading file %s\n", fileName)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ip, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			fmt.Printf("parse ip line error: %v - line: %v", err.Error(), line)
			continue
		}
		if ip.To4() != nil {
			v4 = append(v4, ipnet)
			continue
		}
		if len(ip) == net.IPv6len {
			v6 = append(v6, ipnet)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(v4) == 0 && len(v6) == 0 {
		return nil, fmt.Errorf("cidr list is empty, data file: %v", fileName)
	}
	l.v4 = v4
	l.v6 = v6

	fmt.Printf("loaded file: %d ipnet parsed\n", len(v4)+len(v6))
	return l, nil
}

func (l *CIDRList) contains(ipString string) (bool, error) {
	ipstr := strings.TrimSpace(ipString)
	if ipstr == "" {
		return false, fmt.Errorf("input is empty")
	}
	ip := net.ParseIP(ipstr)
	if len(ip.To4()) == net.IPv4len {
		for _, cidr := range l.v4 {
			if cidr.Contains(ip) {
				return true, nil
			}
		}
	}
	if len(ip) == net.IPv6len {
		for _, cidr := range l.v6 {
			if cidr.Contains(ip) {
				return true, nil
			}
		}
	}
	return false, nil
}

type IpChecker struct {
	l *CIDRList
}

func CreateIpChecker(fileName string) (*IpChecker, error) {
	ch := new(IpChecker)
	l, err := createCIDRList(fileName)
	if err != nil {
		return nil, err
	}
	ch.l = l
	return ch, nil
}

func (ch *IpChecker) RequestHandler(w http.ResponseWriter, r *http.Request) {
	ipstr := r.FormValue("ip")
	contain, err := ch.l.contains(ipstr)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"contains": contain})
}
