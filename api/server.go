package api

import (
	"erli.ng/gotab/storage"
	"github.com/gin-gonic/gin"
)

func CreateServer(disk storage.Disk) *gin.Engine {
	r := gin.Default()

	r.GET("/v1/:partition/:entity", func(ctx *gin.Context) {
		partitionName := ctx.Param("partition")
		entityName := ctx.Param("entity")

		partition := disk.GetPartition(partitionName)
		entity, err := partition.GetEntity(entityName)

		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, entity.EntityData)
	})

	r.POST("/v1/:partition/:entity", func(ctx *gin.Context) {
		partitionName := ctx.Param("partition")
		entityName := ctx.Param("entity")
		data := &storage.EntityData{}
		if err := ctx.ShouldBind(&data); err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		entity := storage.Entity{
			DiskName:      disk.Name,
			PartitionName: partitionName,
			Name:          entityName,
			EntityData:    *data,
		}

		partition := disk.GetPartition(partitionName)
		err := partition.UpdateEntity(entity)

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Status(200)
	})

	r.DELETE("/v1/:partition/:entity", func(ctx *gin.Context) {
		partitionName := ctx.Param("partition")
		entityName := ctx.Param("entity")

		partition := disk.GetPartition(partitionName)
		err := partition.DeleteEntity(storage.Entity{
			DiskName:      disk.Name,
			PartitionName: partitionName,
			Name:          entityName,
		})

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Status(200)
	})

	r.DELETE("/v1/:partition", func(ctx *gin.Context) {
		partitionName := ctx.Param("partition")

		partition := disk.GetPartition(partitionName)
		err := partition.Delete()

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Status(200)
	})

	r.GET("/v1/:partition", func(ctx *gin.Context) {
		partitionName := ctx.Param("partition")

		partition := disk.GetPartition(partitionName)
		entities, err := partition.ListEntities()

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, entities)
	})

	r.GET("/v1", func(ctx *gin.Context) {
		partitions, err := disk.ListPartitionNames()

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, partitions)
	})

	r.GET("/", func(ctx *gin.Context) {

		routes := make(map[string][]string)

		for _, v := range r.Routes() {
			methods := routes[v.Path]
			if methods == nil {
				methods = make([]string, 0)
			}
			methods = append(methods, v.Method)
			routes[v.Path] = methods
		}

		ctx.JSON(200, routes)
	})

	return r
}
