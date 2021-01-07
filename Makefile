all:
	go run servidor_dns/servidor_dns.go

firewall:
	sudo systemctl stop firewalld

clean:
    rm servidor_dns/zf_files/*.zf
    rm servidor_dns/zf_files/*.log  

.PHONY: cliente
.PHONY: admin

cliente:
    go run cliente/cliente.go

admin:
    go run admin/admin.go 