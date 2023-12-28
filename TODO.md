Basic go project

1. Create a new go project that prints hello world
1. Create a script to cross compile that project for
   - ubuntu arm64
   - ubuntu amd64
   - windows amd64
1. read in a yaml config file and print it out
1. map said yaml file to a struct thats kindof what we want and print it

Connect said project to nats

1. add the NATS package and connect to a NATS server, logging out the messages from a fake queue
1. post a message to said queue
1. when we receive a message, spawn a process that does something
1. implement ping pong with the NATS request/response topology

Implement basic shutdown

1. Use ansible to move the binary to the utility machine for testing
1. implement basic shutdown and test it on the util machine

---

1. add to the gitlab server (check if GL has any advanced shutdown things that need doing and do them if so)

Connect the UPS watchers to the system

1. install the nats cli on the dns machines
1. connect NUTs to the system by having it publish a message
1. Test that that message shuts down the util and gitlab servers

Roll it out to most clients

1. Roll this out to k8s
1. Roll this out to j5-box
1. Roll this out to the office and LR PCs

Get the nas and pi switch to shut down

1. Figure out the architecture to tell when all the pis are shut down
   - Should HASS and NAS be responsible for counting or should something else?
   - Can we run arbitrary scripts on truenas?
   - Does HASS have an API for what we want?
   - What is the default state of the kasa switches? When power comes on, will they stay off or turn on?
1. Implement the shutdown of the NAS and the pi switch

Advanced use cases

1. Rearrange the rack as described in the other doc

later

- can we run this on truenas or do we need to use their UPS tooling?
