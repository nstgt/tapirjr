# Playground

Running tapirjr with GoBGP for instant test in docker-compose environment.

## Preparation

* You need to install GoBGP (v2.0.0) to your host machine. See [GoBGP docs](https://github.com/osrg/gobgp).

## Case.1 Development Set

Instant environmet for developing tapirjr.<br>
Run following containers.
* rsN: GoBGP route server container which peer with clients
* cliN: GoBGP client container

#### Start

```
% cd development
% docker-compose up -d
Creating network "development_br01" with driver "bridge"
Creating rs2  ... done
Creating cli2 ... done
Creating cli3 ... done
Creating cli1 ... done
Creating rs3  ... done
Creating rs1  ... done
% docker ps -a
CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS              PORTS                      NAMES
66170a879fb3        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50022->50022/tcp   cli2
d98a5048f020        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50011->50011/tcp   cli1
0b220f49d7e6        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50020->50020/tcp   rs2
5839960a05c5        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50033->50033/tcp   cli3
077ab79b4586        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50010->50010/tcp   rs1
9a81b0e1d673        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   20 seconds ago      Up 18 seconds       0.0.0.0:50030->50030/tcp   rs3
% gobgp -u 127.0.0.1 -p 50011 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.10      65534 00:10:04 Establ      |        0         0
2001:db8:1::10 65534 00:10:05 Establ      |        0         0
% gobgp -u 127.0.0.1 -p 50022 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.20      65534 00:10:05 Establ      |        1         1
2001:db8:1::20 65534 00:10:09 Establ      |        0         0
% gobgp -u 127.0.0.1 -p 50033 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.30      65534 00:10:10 Establ      |        0         0
2001:db8:1::30 65534 00:10:11 Establ      |        0         0
```

#### Play

```
% cd cmd/tapirjr
% go install

## start tapirjr as sender for rs1
% tapirjr sender -g 127.0.0.1:50010 -s 127.0.0.1:50052
2019/01/21 08:13:31 sendder start...

## start tapirjr as receiver for rs2
% tapirjr receiver -g 127.0.0.1:50020 -p 127.0.0.1:50052
2019/01/21 08:13:28 receiver start...

## add prefix to cli1's rib
% gobgp -u 127.0.0.1 -p 50011 global rib add 1.1.0.0/24 origin igp nexthop 10.0.1.11

## then it is advertised to rs1
% gobgp -u 127.0.0.1 -p 50010 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65011                00:07:10   [{Origin: i}]

## from rs1 to rs2, the prefix is sent via tapirjr
% gobgp -u 127.0.0.1 -p 50020 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65011                00:10:36   [{Origin: i}]

## finally this prefix is sent to another client, cli2
% gobgp -u 127.0.0.1 -p 50022 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65534 65011          00:02:22   [{Origin: i}]
```

#### Clean Up

```
% docker-compose down
```

## Case.2 Full Set

Run following containers.
* rsN: GoBGP route server container which peer with clients
* cliN: GoBGP client container
* sndN: tapirjr sender container
* rcvN: tapirjr receiver container

#### Start

```
% cd fullset
% docker-compose up -d
Creating network "fullset_br01" with driver "bridge"
Creating network "fullset_br02" with driver "bridge"
Creating cli1 ... done
Creating rs3  ... done
Creating rcv3 ... done
Creating cli3 ... done
Creating rs1  ... done
Creating cli2 ... done
Creating rs2  ... done
Creating rcv1 ... done
Creating rcv2 ... done
Creating snd1 ... done
Creating snd3 ... done
Creating snd2 ... done
% docker ps -a
CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS              PORTS                      NAMES
e5f4be1194f9        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   23 seconds ago      Up 22 seconds                                  snd2
d0106ff040b4        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   23 seconds ago      Up 22 seconds                                  snd3
adeba56927e2        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   23 seconds ago      Up 22 seconds                                  snd1
8e726cc9ad9f        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 22 seconds       0.0.0.0:50020->50020/tcp   rs2
12721f7855ae        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   27 seconds ago      Up 23 seconds       0.0.0.0:50052->50052/tcp   rcv2
33b4273b7588        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   27 seconds ago      Up 23 seconds       0.0.0.0:50051->50051/tcp   rcv1
af8e5eb4251e        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 23 seconds       0.0.0.0:50022->50022/tcp   cli2
d9196a963f12        nstgt/tapirjr:v0.0.0        "/usr/local/bin/tapi…"   27 seconds ago      Up 23 seconds       0.0.0.0:50053->50053/tcp   rcv3
da6e3d7c6631        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 23 seconds       0.0.0.0:50010->50010/tcp   rs1
2ee1c2ca5447        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 23 seconds       0.0.0.0:50011->50011/tcp   cli1
126c80b66cf0        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 23 seconds       0.0.0.0:50030->50030/tcp   rs3
698162409c19        nstgt/docker-gobgp:latest   "/usr/local/bin/gobg…"   27 seconds ago      Up 23 seconds       0.0.0.0:50033->50033/tcp   cli3
% gobgp -u 127.0.0.1 -p 50011 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.10      65534 00:05:17 Establ      |        0         0
2001:db8:1::10 65534 00:05:15 Establ      |        0         0
% gobgp -u 127.0.0.1 -p 50022 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.20      65534 00:05:23 Establ      |        0         0
2001:db8:1::20 65534 00:05:19 Establ      |        0         0
% gobgp -u 127.0.0.1 -p 50033 nei
Peer              AS  Up/Down State       |#Received  Accepted
10.0.1.30      65534 00:05:29 Establ      |        0         0
2001:db8:1::30 65534 00:05:23 Establ      |        0         0
```


#### Play

```
## add prefix to cli1's rib
% gobgp -u 127.0.0.1 -p 50011 global rib add 1.1.0.0/24 origin igp nexthop 10.0.1.11

## then it is advertised to rs1
% gobgp -u 127.0.0.1 -p 50010 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65011                00:07:10   [{Origin: i}]

## from rs1 to rs2, the prefix is sent via tapirjr
% gobgp -u 127.0.0.1 -p 50020 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65011                00:10:36   [{Origin: i}]

## finally this prefix is sent to another client, cli2
% gobgp -u 127.0.0.1 -p 50022 global rib
   Network              Next Hop             AS_PATH              Age        Attrs
*> 1.1.0.0/24           10.0.1.11            65534 65011          00:02:22   [{Origin: i}]
```

#### Clear Up

```
% docker-compose down
```