# NetTwerking
How to create a distributed system 101

## Docker START
1. Make sure docker is installed.
2. Go to root directory of project.
3. Run the script "startdocker.sh" (in the script, alter the amount of nodes to be started if you like).
4. The images will now be built, and the specified amount of containers will be started.
5. You will automatically enter an interactive container. Do as you please!
6. When exiting the container, all containers will be closed.

### Docker debugging commands
* docker logs 'id' : Shows all stuff printed to stdout in the container.
