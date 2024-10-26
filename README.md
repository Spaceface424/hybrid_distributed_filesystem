# CS425_MP2_G14

UIUC CS 425 MP2 Repository 

Group Number 14

Ryan Wu (rwu29), Oliver Rogalski (oliverr3)

To run the MP, clone this repository into all of the virtual machines that will be used. Run swim/main.go on the introducer process then provide the introducer IP to any joining processes.

Example:
```
go run swim/main.go [-i introducer_ip] [RPC port] [UDP port]
```

To run commands, open a new terminal and run client/main.go in the following format:
```
go run client/main.go -o host_file -m [command] [input_argument]
```
To run commands on one member, pass the ip address and RPC port into the input_argument.
Example:
```
go run client/main.go -o host_file -m list_mem fa24-cs425-1401.cs.illinois.edu:8000
```

To change the drop rate across all members, pass the desired drop rate into input_argument.
Example:
```
go run client/main.go -o host_file -m drop_rate 0.15
```

Do not pass input_arguments when running enable/disable_sus and list_sus commands. The command will be sent to all VMs.
Example:
```
go run client/main.go -o host_file -m enable_sus
```

Commands:
-  list_mem - Prints the membership list of the member
- list_self - Prints the member's own ID
- leave - Makes the member exit and leave the network
- enable/disable_sus - Enables or disables the suspicion protocol on all machines
- status_sus - Prints whether the suspicion protocl is enabled on a member
- drop_rate - Sets all machines' drop rate to the input_argument
- list_sus - Prints all suspicious members in each machine's membership list
