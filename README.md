# NetTwerking
How to create a distributed system 101

## Docker START
1. Make sure all docker related stuff is installed (docker, docker-compose).
2. Go to the directory where your docker-compose file is.
3. Run command: docker swarm init
4. Run command: docker stack deploy --compose-file docker-compose.yml name-of-stack
5. To enter the interactive node started, look up its id using docker ps.
6. Enter the node by running: docker attach 'id' (Previously stuff printed to stdout will not be shown.)
7. Do what you want in the container.
8. When quitting, run: docker stack rm name-of-stack & docker swarm leave --force

### Docker debugging commands
* docker logs 'id' : Shows all stuff printed to stdout in the container.
* docker stack services 'stack-name' : Shows how many replicas you have up and running of each service.
