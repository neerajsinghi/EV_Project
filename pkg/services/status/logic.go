package status

import (
	"context"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func Logic(startDate, endDate time.Time, city, service string) (interface{}, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -365)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	pipeline := getPipeline(startDate, endDate)

	if city != "" {
		pipeline = getPipelineCity(startDate, endDate, city)
	}
	cursor, err := trestCommon.Aggregate(pipeline, "users")
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var result []interface{}
	for cursor.Next(context.TODO()) {
		var profile interface{}
		if err = cursor.Decode(&profile); err != nil {
			trestCommon.ECLog3(
				"Find profiles",
				err,
				logrus.Fields{
					"filter":   bson.M{},
					"error at": cursor.RemainingBatchLength(),
				})
			return nil, nil
		}
		result = append(result, profile)
	}
	return result, nil
}

func getPipeline(startDate time.Time, endDate time.Time) bson.A {
	pipeline := bson.A{

		bson.D{
			{Key: "$facet", Value: bson.D{
				{Key: "numberOfUsers", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}},
					}},
					bson.D{{Key: "$count", Value: "count"}}}},
				{Key: "idVerified", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "dl_verified", Value: true}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "dlVerified", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "dl_verified", Value: true}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "unverifiedUsers", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "$and", Value: bson.A{
						bson.D{{Key: "dl_verified", Value: false}},
						bson.D{{Key: "id_verified", Value: false}},
					}}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "carbonEmissionSaved", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "total", Value: bson.D{{Key: "$sum", Value: "$carbon_saved"}}},
					}}},
				}},
			}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "station"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
						},
						},
						},
						bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "publicStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "public", Value: true}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
					{Key: "activeStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "status", Value: "available"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "stationStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "charger"},
			{Key: "pipeline", Value: bson.A{bson.D{{Key: "$count", Value: "count"}}}},
			{Key: "as", Value: "chargerStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "booking"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalRides", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
						bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "totalDistance", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
						bson.D{{Key: "$group", Value: bson.D{
							{Key: "_id", Value: nil},
							{Key: "total", Value: bson.D{{Key: "$sum", Value: "$total_distance"}}},
						}}},
						bson.D{{Key: "$project", Value: bson.D{
							{Key: "total", Value: bson.D{{Key: "$round", Value: bson.D{{Key: "$divide", Value: bson.A{"$total", 1000}}}}}},
						}}},
					}},
					{Key: "completedRides", Value: bson.A{

						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "status", Value: "completed"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "bookingStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "iotBike"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalVehicles", Value: bson.A{bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "activeVehicles", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "online"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
					{Key: "vehiclesOnRoad", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "ignition", Value: "true"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "bikeStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "wallet"},
			{Key: "let", Value: bson.D{}},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "totalEarning", Value: bson.D{{Key: "$sum", Value: "$used_money"}}},
					{Key: "totalValueInWallet", Value: bson.D{{Key: "$sum", Value: "$deposited_money"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "totalEarning", Value: bson.D{{Key: "$add", Value: bson.A{"$totalEarning", "$totalValueInWallet"}}}},
					{Key: "totalValueInWallet", Value: 1},
					{Key: "_id", Value: 0},
				}}},
			}},
			{Key: "as", Value: "walletStats"},
		}}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "numberOfUsers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$numberOfUsers.count", 0}}}},
			{Key: "idVerified", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$idVerified.count", 0}}}},
			{Key: "dlVerified", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$dlVerified.count", 0}}}},
			{Key: "unverifiedUsers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$unverifiedUsers.count", 0}}}},
			{Key: "totalStations", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.totalStations.count", 0}}}},
			{Key: "totalPublicStations", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.publicStations.count", 0}}}},
			{Key: "totalActiveStation", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.activeStations.count", 0}}}},
			{Key: "totalCo2Emission", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.carbonEmissionSaved.total", 0}}}},
			{Key: "totalChargers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$chargerStats.count", 0}}}},
			{Key: "totalRides", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.totalRides.count", 0}}}},
			{Key: "totalDistance", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.totalDistance.total", 0}}}},
			{Key: "totalCompletedRides", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.completedRides.count", 0}}}},
			{Key: "totalVehicles", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.totalVehicles.count", 0}}}},
			{Key: "totalActiveVeficles", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.activeVehicles.count", 0}}}},
			{Key: "totalVehicleOnRoad", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.vehiclesOnRoad.count", 0}}}},
			{Key: "totalEarning", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$walletStats.totalEarning", 0}}}},
			{Key: "totalValueInWallet", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$walletStats.totalValueInWallet", 0}}}},
		}}},
	}
	return pipeline
}

func getPipelineCity(startDate time.Time, endDate time.Time, city string) bson.A {
	pipeline := bson.A{

		bson.D{
			{Key: "$facet", Value: bson.D{
				{Key: "numberOfUsers", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}},
					}},
					bson.D{{Key: "$count", Value: "count"}}}},
				{Key: "idVerified", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "dl_verified", Value: true}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "dlVerified", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "dl_verified", Value: true}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "unverifiedUsers", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "$and", Value: bson.A{
						bson.D{{Key: "dl_verified", Value: false}},
						bson.D{{Key: "id_verified", Value: false}},
					}}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "carbonEmissionSaved", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "total", Value: bson.D{{Key: "$sum", Value: "$carbon_saved"}}},
					}}},
				}},
			}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "station"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "city", Value: city},
						},
						},
						},
						bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "publicStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "city", Value: city},
							{Key: "public", Value: true}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
					{Key: "activeStations", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "city", Value: city},

							{Key: "status", Value: "available"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "stationStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "charger"},
			{Key: "pipeline", Value: bson.A{bson.D{{Key: "$count", Value: "count"}}}},
			{Key: "as", Value: "chargerStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "booking"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalRides", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "city", Value: city},
						}}},

						bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "totalDistance", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}, {Key: "city", Value: city},
						}}},
						bson.D{{Key: "$group", Value: bson.D{
							{Key: "_id", Value: nil},
							{Key: "total", Value: bson.D{{Key: "$sum", Value: "$total_distance"}}},
						}}},
						bson.D{{Key: "$project", Value: bson.D{
							{Key: "total", Value: bson.D{{Key: "$round", Value: bson.D{{Key: "$divide", Value: bson.A{"$total", 1000}}}}}},
						}}},
					}},
					{Key: "completedRides", Value: bson.A{

						bson.D{{Key: "$match", Value: bson.D{
							{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
							{Key: "city", Value: city},

							{Key: "status", Value: "completed"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "bookingStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "iotBike"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$facet", Value: bson.D{
					{Key: "totalVehicles", Value: bson.A{bson.D{{Key: "$count", Value: "count"}}}},
					{Key: "activeVehicles", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "online"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
					{Key: "vehiclesOnRoad", Value: bson.A{
						bson.D{{Key: "$match", Value: bson.D{{Key: "ignition", Value: "true"}}}},
						bson.D{{Key: "$count", Value: "count"}},
					}},
				}}},
			}},
			{Key: "as", Value: "bikeStats"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "wallet"},
			{Key: "let", Value: bson.D{}},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "totalEarning", Value: bson.D{{Key: "$sum", Value: "$used_money"}}},
					{Key: "totalValueInWallet", Value: bson.D{{Key: "$sum", Value: "$deposited_money"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "totalEarning", Value: bson.D{{Key: "$add", Value: bson.A{"$totalEarning", "$totalValueInWallet"}}}},
					{Key: "totalValueInWallet", Value: 1},
					{Key: "_id", Value: 0},
				}}},
			}},
			{Key: "as", Value: "walletStats"},
		}}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "numberOfUsers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$numberOfUsers.count", 0}}}},
			{Key: "idVerified", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$idVerified.count", 0}}}},
			{Key: "dlVerified", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$dlVerified.count", 0}}}},
			{Key: "unverifiedUsers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$unverifiedUsers.count", 0}}}},
			{Key: "totalStations", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.totalStations.count", 0}}}},
			{Key: "totalPublicStations", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.publicStations.count", 0}}}},
			{Key: "totalActiveStation", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.activeStations.count", 0}}}},
			{Key: "totalCo2Emission", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$stationStats.carbonEmissionSaved.total", 0}}}},
			{Key: "totalChargers", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$chargerStats.count", 0}}}},
			{Key: "totalRides", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.totalRides.count", 0}}}},
			{Key: "totalDistance", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.totalDistance.total", 0}}}},
			{Key: "totalCompletedRides", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bookingStats.completedRides.count", 0}}}},
			{Key: "totalVehicles", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.totalVehicles.count", 0}}}},
			{Key: "totalActiveVeficles", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.activeVehicles.count", 0}}}},
			{Key: "totalVehicleOnRoad", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$bikeStats.vehiclesOnRoad.count", 0}}}},
			{Key: "totalEarning", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$walletStats.totalEarning", 0}}}},
			{Key: "totalValueInWallet", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$walletStats.totalValueInWallet", 0}}}},
		}}},
	}
	return pipeline
}
