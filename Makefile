all:
	#go run servidor_dns/servidor_dns.go 1

firewall:
	sudo systemctl start firewalld

clean:
	rm servidor_dns/zf_files/*.zf
	rm servidor_dns/zf_files/*.log	

#cliente:
	#go run cliente/cliente.go 

#broker:
	#go run broker/broker_server.go 

#admin:
	#go run admin/admin.go 