package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

func serverDNS() {
	addr := net.UDPAddr{
		Port: config.DnsPort,
		IP:   net.ParseIP(config.DnsListenAddress),
	}
	u, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
		return
	}

	log.Println("listen udp ", addr.String())
	// Wait to get request on that port
	for {
		tmp := make([]byte, 1024)
		_, addr, _ := u.ReadFrom(tmp)
		clientAddr := addr
		//logs.Println(addr)
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		serveDNS(u, clientAddr, tcp)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	//var ip string
	var err error
	log.Println(clientAddr, "query:", string(request.Questions[0].Name))
	dns2 := redisServer.ReadDnsByDomain(string(request.Questions[0].Name))
	if dns2.RecodeType == "a" {
		//log.Println("aaaaaaa")
		dnsAnswer.Type = layers.DNSTypeA
		dnsAnswer.IP = net.ParseIP(dns2.RecodeValue).To4()
		dnsAnswer.Name = []byte(request.Questions[0].Name)
		//fmt.Println(request.Questions[0].Name)
		dnsAnswer.Class = layers.DNSClassIN
		replyMess.QR = true
		replyMess.ANCount = 1
		replyMess.OpCode = layers.DNSOpCodeNotify
		replyMess.AA = true
		replyMess.Answers = append(replyMess.Answers, dnsAnswer)
		replyMess.ResponseCode = layers.DNSResponseCodeNoErr
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = replyMess.SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		_, _ = u.WriteTo(buf.Bytes(), clientAddr)
	} else if dns2.RecodeType == "aaaa" {
		//a, _, _ := net.ParseCIDR(ip + "/24")
		dnsAnswer.Type = layers.DNSTypeAAAA
		dnsAnswer.IP = net.ParseIP(dns2.RecodeValue).To16()
		dnsAnswer.Name = []byte(request.Questions[0].Name)
		//fmt.Println(request.Questions[0].Name)
		dnsAnswer.Class = layers.DNSClassIN
		replyMess.QR = true
		replyMess.ANCount = 1
		replyMess.OpCode = layers.DNSOpCodeNotify
		replyMess.AA = true
		replyMess.Answers = append(replyMess.Answers, dnsAnswer)
		replyMess.ResponseCode = layers.DNSResponseCodeNoErr
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = replyMess.SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		_, _ = u.WriteTo(buf.Bytes(), clientAddr)
	} else if dns2.RecodeType == "cname" {
		//a, _, _ := net.ParseCIDR(ip + "/24")
		dnsAnswer.Type = layers.DNSTypeCNAME
		dnsAnswer.CNAME = []byte(dns2.RecodeValue)
		dnsAnswer.Name = []byte(request.Questions[0].Name)
		//fmt.Println(request.Questions[0].Name)
		dnsAnswer.Class = layers.DNSClassIN
		replyMess.QR = true
		replyMess.ANCount = 1
		replyMess.OpCode = layers.DNSOpCodeNotify
		replyMess.AA = true
		replyMess.Answers = append(replyMess.Answers, dnsAnswer)
		replyMess.ResponseCode = layers.DNSResponseCodeNoErr
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = replyMess.SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		_, _ = u.WriteTo(buf.Bytes(), clientAddr)
	} else {
		//a, _, _ := net.ParseCIDR(ip + "/24")
		dnsAnswer.Type = layers.DNSTypeA
		dnsAnswer.IP = net.ParseIP(dns2.RecodeValue)
		dnsAnswer.Name = []byte(request.Questions[0].Name)
		fmt.Println(request.Questions[0].Name)
		dnsAnswer.Class = layers.DNSClassIN
		replyMess.QR = true
		replyMess.ANCount = 1
		replyMess.OpCode = layers.DNSOpCodeNotify
		replyMess.AA = true
		replyMess.Answers = append(replyMess.Answers, dnsAnswer)
		replyMess.ResponseCode = layers.DNSResponseCodeNXDomain
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
		err = replyMess.SerializeTo(buf, opts)
		if err != nil {
			panic(err)
		}
		_, _ = u.WriteTo(buf.Bytes(), clientAddr)
	}

}
