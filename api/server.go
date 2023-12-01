package api

import (
	"erli.ng/gotab/storage"
	"github.com/gin-gonic/gin"
)

func CreateServer(disk storage.Disk) *gin.Engine {
	r := gin.Default()

	r.GET("/disk/:partition/:entity", func(ctx *gin.Context) {
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

	r.POST("/disk/:partition/:entity", func(ctx *gin.Context) {
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

		ctx.JSON(200, entity.EntityData)
	})

	r.DELETE("/disk/:partition/:entity", func(ctx *gin.Context) {
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

	r.DELETE("/disk/:partition", func(ctx *gin.Context) {
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

	r.GET("/disk/:partition", func(ctx *gin.Context) {
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

	r.GET("/disk", func(ctx *gin.Context) {
		partitions, err := disk.ListPartitionNames()

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, partitions)
	})
	return r
}
