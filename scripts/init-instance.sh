NM=$1

multipass launch --name ${NM} jammy --memory 2G --disk 8G --cpus 2

multipass transfer install-docker.sh ${NM}:/home/ubuntu/install-docker.sh
multipass exec ${NM} -- sh -x /home/ubuntu/install-docker.sh
