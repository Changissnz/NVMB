#GO_GET(go_redis github.com/hoisie/redis)
GO_GET(gonumGet gonum.org/v1/gonum/) 
GO_GET(flsGet github.com/stoicperlman/fls)

GO_GET(km github.com/biogo/cluster/kmeans)

ADD_GO_INSTALLABLE_PROGRAM(golem
                           main.go
                           gonumGet
                           flsGet 
                           km)