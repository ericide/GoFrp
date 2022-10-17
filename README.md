# GoFrp

## Introduce

A fast reverse proxy to help you expose a local server behind a NAT or firewall to the internet

 - Save ports, only one port needs to be opened on the server side
 - Full platform support
 - Easy. Simple to use and configuration. No need for docker, it runs directly on your machine with minimal memory consumption and cpu consumption, and is also very power efficient

## Example

```
//on machine A which has a public IP (a.b.c.d)
./frp -m server -h 0.0.0.0 -p 10000

//on machine B which behind a NAT 
./frp -m client -h a.b.c.d -p 10000 -lh localhost -lp 443

// Now you can visit machine B`s 443 port by a.b.c.d:10000
```

## How to use
Command
```
  -h string
        server host (default "0.0.0.0")
  -p int
        server bind port (default 10000)
  -lh string
        local host (default "localhost")
  -lp int
        local bind port (default 443)
  -m string
        run mode (default "server")
  -pwd string
        password for connect (default "12345678")
```

## Completion

- [x] Encryption signalling channel
- [x] Use sigle port handling data transmit and signalling transmit
- [ ] UDP Method


## Some more to say
                       
Real Client <------> Frp Server(A) <======> Frp Client(B) <------> Real Server

In this situation, real server dosn`t have a public IP.

When you're transferring information between A and B. Only the signalling information is encrypted to ensure that the current status of the program is not known to the outside world, and to reject interfering links from outside. This ensures a reliable link between the two ends of the software.

This software uses tcp as a linking scheme. It exposes a local port to the outside. So information from your local server will be transmitted to the outside via the network connection, which can also lead to all kinds of attacks, information interception and even man-in-the-middle attacks. So your service has to take this into account and take certain precautions. For example, if you set up an http service, you should use https for the protection of your information.

You may ask why only the signalling information is encrypted and not the transmitted information. The purpose of encrypting the signalling information is to ensure that the Frp pair can keep their working state from being exposed to the outside network, and to give the Frp pair a means of confirming each other's identity when a new link is created, preventing interference by other attackers. (If Frp is not encrypted and signalled carefully, this could result in Frp being disrupted by an attack, but of course if your application has its own encryption, this does not affect its security, only the link availability). Even if Frp encrypts the actual linked data messages, eventually these messages will be decrypted and flow out of the Frp system (Frp must ensure that the data entering one end of Frp will flow out at the other end unchanged) into the public network, so encrypting communication messages in the Frp bipartite does not make sense, but rather reduces the efficiency of communication and increases the energy consumption of the software as well as the memory and CPU usage.