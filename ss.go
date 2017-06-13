package main

import (
	"bufio"

	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"strings"
)

func privateKeyPath() string {
	return os.Getenv("HOME") + "/gitkeys/id_rsa"
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	//fmt.Println(keyPath)
	buff := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAttLRDXLQyKWpYE9ZvNwKjPvLXOKd9zIO78QTgykJTM9jJ095
CnCDilBXufnATIyqxf14/WjYC3DxMDFpq25FSGH+Uwo65y42Tx/hJeYPO7GTUsDa
uGtDvJC6yCzumsE5ZhqS/8+KCENbt8IqDKHMmbHGLQDMvONbHV/K5clMnjFsuMI1
+kYeJFZ33hCBamthjkh6TW2H8twx1CXiYPihVtF/DaWRI9E+h6kgKsvgYyNActtq
l0RKqLUXRKdUZBpmS9wVEZjF1H0MEkX42tc+8gVcYOao7jPTC6/CGo1wGlpey7w2
K6q9wNeyBGfaeryDbxLwLVoAAdqdnCxlc+/IvwIDAQABAoIBAATGnQ9s4GddbH9H
k2CLnvrv2PAkO1wpwpog1SGoAMYw5LYyLUdIiScj0ibZj1xwkEV7yZ1VD8+8s5A+
ujPwPX8WkD47Fi7T1jda7da/m8ZDbUXMI+qNYseaQEbaZaFFSuqP5ycnlMOCfvLd
75tE3sNlEXg7fR2yEr9BsVsvVwEK3i2Mcz275zMMUg9XjLhpKsHzAISrDHL3yMCj
9Jhs8nHF+HkIOVde3Ly1dgd83BVuhzrPeJBFQuQ3BOFQ4SS1rMQdK+gH3I3tegU6
NgcJHKnS80KVUzkTUVbRTqt7UWZ8F6Q93I2gbrFxuVie633h4L0lSBPw5Ax6/LOJ
91SyHoECgYEA3nAGgAqjB/TblOXAntNmKCNYSuPUFm6f9Evnc+zPdAx0m9V0rse0
XqAHcANt6jV89z5jchLHgAYicpC6Wriu4SlLxakGqsxdwPFefUw4BixTEubdZe/g
N3OAIZ0Eaoo5QsMpJegd0wgvq0sNnoWLTkTBZumXFEXHvDkhNFXOKvsCgYEA0mij
qtfW/AZH4Tndof9pPIpgQ8xbEhy5SqCpT+Rg+Nin6ovTOI0hC6HrcmeIfjWiAAhm
t/VC5e2m0OoXrEXqGmZ1/UjnjQHWVEtWP47AKfpjxQ97PA4YoUGo0zaolzzEFM5U
vTNA6AxVESPHrOJyAxF+NooyLLnTTUVWIJHzrg0CgYB0GQhrgCHDn1uUha5Zt4DU
Zk5JGEy0QJ0gBxYQ/YLx0SZzx5+VMgrEcMYxArk1yyEkct24xnB2M717CmsZutcc
Ek/IJQaj0vMEJ3bn8wYywqPBc9oOwHrItnIkGS4a5XYpkG9Dp7kZUmZ/Azdii9U/
zscbDcSbAijT5wWbqUVoTQKBgQCoeDQjxLJUFMtU4Lo+zXx7huhRIL0CoZES2dT3
LQsf9IluWQqESyvcXodgkNlPBK5zjEaCoJQx+bkJqYXO4CPzg3qRlOAhnQj5cWDb
fvcKJXvg+uZXTYoXA7WjeC5A+dyeNB7RZspfghBSqu1j1eQn5MfD758BBMDVK+Es
LvCHuQKBgH+g8XWVVlwmz+4MIqZtwOagrfWM+cS88HEvVnTSMW8s1Lnw8KA4Z/uP
e+0m/pXJEsluY0RXfNxFBLnrP9wUaGCwhkvq59BFj3NeDEeormEmfWNvbvofaQZf
37IkOK+UyMxyfwV4sx29mBuF1p3d4EZrq/zYEE4akKEVDj5khLq/
-----END RSA PRIVATE KEY-----
`)

	return ssh.ParsePrivateKey(buff)

}

func main() {

	server := "10.4.3.30:22"
	user := "root"

	key, err := parsePrivateKey(privateKeyPath())

	if err != nil {
		panic("Error: " + err.Error())
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	in, _ := session.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	// Accepting commands
	time.Sleep(3 * time.Second)

	str := "nc localhost 8089\r\n"
	str2 := "this is a dummy message\r\n"
	fmt.Fprint(in, str)
	fmt.Fprintf(in, str2)

	for str != "exit\n" {
		reader := bufio.NewReader(os.Stdin)
		str, _ = reader.ReadString('\n')
		fmt.Fprint(in, str)

	}

}

func scanConfig() string {
	config, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	config = strings.Trim(config, "\n")
	return config
}
