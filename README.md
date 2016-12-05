# Parallel server architecture for position and movement of elements in a 2D space

For fun project to build a scalable architecture for tracking movement and position of elements in a 2D game. Using multiple positionServers to store and pass elements between each other and a positionServerSupervisor to lockstep them all, take a look at serverDiagram.png for a simplistic view. WIP.

## RUN
Use 'go build' and then run executable, -h for flag options. Not much to see without a client renderer (coming), but use the verbose option to show logs of what's going on.

## TEST    
'go test
At the moment ~50% coverage, not great!

## TODO
Metric Server
- watch global entity number
- watch ticks/second

Client Renderer
- SDL window to see the elements moving on the screen
- can use my own Go game engine
