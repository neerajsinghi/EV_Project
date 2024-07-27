package status

import (
	"bikeRental/pkg/repo/bikeDevice"
	iotbike "bikeRental/pkg/repo/iot_bike"
	"bikeRental/pkg/services/motog"
	"context"
	"errors"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

/*
[{

	  "deviceId": 310823018,
	  "alarm": "",
	  "batteryLevel": 49.2,
	  "course": "",
	  "dealer": "",
	  "deviceFixTime": "",
	  "deviceImei": "",
	  "harshAccelerationHistory": null,
	  "harshBrakingHistory": null,
	  "ignition": "0",
	  "lastUpdate": "2024-07-27 15:05:45",
	  "latitude": "",
	  "location": {
	    "type": "Point",
	    "coordinates": [
	      77.200415,
	      28.582604
	    ]
	  },
	  "longitude": "",
	  "name": "MVU310823018",
	  "phone": "",
	  "posId": 0,
	  "speed": "0",
	  "status": "",
	  "totalDistance": "286",
	  "type": "moto",
	  "valid": 0,
	  "booking": {
	    "_id": {
	      "$oid": "66a507ad0b32d51db4b27d98"
	    },
	    "profile_id": "668e37fa55b793c4a80f9a52",
	    "device_id": 310823018,
	    "start_time": 1722091437,
	    "end_time": 1722092562,
	    "start_km": 281,
	    "end_km": 281,
	    "total_distance": 0,
	    "return": {
	      "location": "",
	      "time": "",
	      "product_images": [
	        "",
	        ""
	      ],
	      "damages": [
	        ""
	      ],
	      "front_picture": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/79FE9FA0-1D25-4698-9086-1CDC936835CA.jpg?alt=media&token=74aa853c-eb10-478f-bb58-6c7ba2138518",
	      "back_picture": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/A3866C9B-F95E-4593-AF9B-CEA44A85E3D0.jpg?alt=media&token=b7511223-9db4-438d-919c-aa14b034b725",
	      "left_picture": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/911AE54E-799F-4672-9715-8FA6E8E859F0.jpg?alt=media&token=bf4a0c11-0c81-462f-b7c1-f8fbcc530830",
	      "right_picture": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/70C579A9-7BE1-4422-A8CF-04E8F05DB883.jpg?alt=media&token=b85805e1-b746-4531-aa8b-422a05877302",
	      "front_desc": "",
	      "back_desc": "",
	      "left_desc": "",
	      "right_desc": ""
	    },
	    "price": 10,
	    "status": "completed",
	    "vehicle_type": "moto",
	    "booking_type": "hourly",
	    "plan": {
	      "_id": {
	        "$oid": "66a1eec1a0b86712f6245b25"
	      },
	      "name": "",
	      "city": "Delhi",
	      "vehicle_type": "Normal",
	      "charger_type": "",
	      "type": "hourly",
	      "description": "",
	      "starting_minutes": 0,
	      "ending_minutes": 30,
	      "every_x_minutes": 0,
	      "price": 10,
	      "deposit": null,
	      "validity": "",
	      "discount": 0,
	      "is_active": true,
	      "status": "",
	      "created_time": {
	        "$date": "2024-07-25T06:20:49.230Z"
	      }
	    },
	    "created_time": {
	      "$date": "2024-07-27T14:43:57.558Z"
	    },
	    "starting_station_id": "66a268f5e768c56691fa10e6",
	    "ending_station_id": "66a268f5e768c56691fa10e6",
	    "carbon_emission_saved": 0,
	    "starting_station": {
	      "_id": {
	        "$oid": "66a268f5e768c56691fa10e6"
	      },
	      "name": "Chanakyapuri",
	      "description": "",
	      "short_name": "Chanakya",
	      "address": {
	        "address": "Madhu Limaye Mar, Chanakyapuri",
	        "country": "India",
	        "pin": "110021",
	        "city": "Delhi",
	        "state": "Delhi"
	      },
	      "location": {
	        "type": "Point",
	        "coordinates": [
	          77.20027850454505,
	          28.582495809040516
	        ]
	      },
	      "active": true,
	      "group": "",
	      "supervisor_id": "",
	      "stock": 0,
	      "public": true,
	      "status": "available",
	      "services_available": [
	        "hourly",
	        "rental"
	      ],
	      "update_at": {
	        "$date": "2024-07-25T16:58:21.786Z"
	      },
	      "created_time": {
	        "$date": "2024-07-25T15:02:13.483Z"
	      },
	      "location_polygon": null
	    },
	    "ending_station": {
	      "_id": {
	        "$oid": "66a268f5e768c56691fa10e6"
	      },
	      "name": "Chanakyapuri",
	      "description": "",
	      "short_name": "Chanakya",
	      "address": {
	        "address": "Madhu Limaye Mar, Chanakyapuri",
	        "country": "India",
	        "pin": "110021",
	        "city": "Delhi",
	        "state": "Delhi"
	      },
	      "location": {
	        "type": "Point",
	        "coordinates": [
	          77.20027850454505,
	          28.582495809040516
	        ]
	      },
	      "active": true,
	      "group": "",
	      "supervisor_id": "",
	      "stock": 0,
	      "public": true,
	      "status": "available",
	      "services_available": [
	        "hourly",
	        "rental"
	      ],
	      "update_at": {
	        "$date": "2024-07-25T16:58:21.786Z"
	      },
	      "created_time": {
	        "$date": "2024-07-25T15:02:13.483Z"
	      },
	      "location_polygon": null
	    },
	    "coupon_code": "",
	    "discount": 0,
	    "green_points": 0,
	    "carbon_saved": 0,
	    "city": "Delhi",
	    "ride_time_remaining": 161,
	    "update_time": {
	      "$date": "2024-07-27T15:02:27.032Z"
	    },
	    "time_booked": 18,
	    "updated_at": {
	      "$date": "2024-07-27T15:02:43.403Z"
	    }
	  },
	  "planId": "66a1eec1a0b86712f6245b25",
	  "userID": {
	    "$oid": "668e37fa55b793c4a80f9a52"
	  },
	  "plan": {
	    "_id": {
	      "$oid": "66a1eec1a0b86712f6245b25"
	    },
	    "name": "",
	    "city": "Delhi",
	    "vehicle_type": "Normal",
	    "charger_type": "",
	    "type": "hourly",
	    "description": "",
	    "starting_minutes": 0,
	    "ending_minutes": 30,
	    "every_x_minutes": 0,
	    "price": 10,
	    "deposit": null,
	    "validity": "",
	    "discount": 0,
	    "is_active": true,
	    "status": "",
	    "created_time": {
	      "$date": "2024-07-25T06:20:49.230Z"
	    }
	  },
	  "planType": "hourly",
	  "price": 10,
	  "profile": {
	    "_id": {
	      "$oid": "668e37fa55b793c4a80f9a52"
	    },
	    "email": "",
	    "status": "created",
	    "status_bool": null,
	    "joining_date": "",
	    "name": "Karan Chauhan",
	    "dob": "9/27/1991",
	    "designation": null,
	    "gender": "Male",
	    "phone_no": "+919999123111",
	    "phone_otp": "328919",
	    "roles": "user",
	    "phone_no_verified": false,
	    "address": null,
	    "about": null,
	    "url_to_profile_image": null,
	    "password": "$2a$05$AalNiH.SDB023WdAjta31uuhYzGFk76cVC/1ffFl2lOp2H4JcRA5W",
	    "created_time": {
	      "$date": "2024-07-10T07:27:54.313Z"
	    },
	    "email_login_otp": null,
	    "otp_code": null,
	    "update_time": {
	      "$date": "2024-07-27T15:02:42.352Z"
	    },
	    "email_sent_time": null,
	    "verification_code": null,
	    "password_reset_code": null,
	    "country_code": null,
	    "password_reset_time": {
	      "$date": {
	        "$numberLong": "-62135596800000"
	      }
	    },
	    "last_login_device_id": null,
	    "last_login_device_name": null,
	    "last_login_location": null,
	    "online": null,
	    "dl_verified": true,
	    "dl_image": "",
	    "id_image": "",
	    "id_verified": true,
	    "plan_id": null,
	    "plan": null,
	    "plan_start_time": 0,
	    "plan_active": null,
	    "user_blocked": null,
	    "plan_end_time": 0,
	    "plan_remaining_time": 0,
	    "referral_code": "vdWR0kU5",
	    "referral_code_used": null,
	    "access": null,
	    "terms_and_condition": true,
	    "allow_promotions": false,
	    "firebase_token": "f2rIpJ-bakYWi6g0yPOU7g:APA91bGKzZwoOS23dhI0B8JjPaWqBbbZY5b7y3y8xXLcZrTb_c41oC_4Z8171iFLT1kmx4pZ94jIpNaP4Gedm1hiQq8hEY8fAotOrw6kMFuQ9bc1McL_cHZX4uIExOfz5WFTzVAeCGTr",
	    "total_balance": 0,
	    "total_rides": 0,
	    "green_points": 0,
	    "carbon_saved": 0,
	    "total_travelled": 0,
	    "referred_by": "",
	    "verified_time": {
	      "$date": "2024-07-27T13:28:20.787Z"
	    },
	    "id_back_image": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/FA9009D7-58F8-447B-8E83-BE84A540DD44.jpg?alt=media&token=accbffe4-86f6-407e-bb4c-fa88e2e1f233",
	    "id_front_image": "https://firebasestorage.googleapis.com/v0/b/futureevdemo.appspot.com/o/8F8AF2E4-B0AC-4BFE-B2FA-659255C864B5.jpg?alt=media&token=7ca5728f-6678-4306-b343-7e78faef183a",
	    "service_type": "hourly"
	  }
	}]
*/
type VehicleData struct {
	DeviceId     int     `bson:"deviceId,omitempty"`
	BatteryLevel float64 `bson:"batteryLevel,omitempty"`
	LastUpdate   string  `bson:"lastUpdate,omitempty"`
	Location     struct {
		Type        string    `bson:"type,omitempty"`
		Coordinates []float64 `bson:"coordinates,omitempty"`
	} `bson:"location,omitempty"`
	Name          string `bson:"name,omitempty"`
	Speed         string `bson:"speed,omitempty"`
	TotalDistance string `bson:"totalDistance,omitempty"`
	Type          string `bson:"type,omitempty"`
	Booking       struct {
		ID                  primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		StartTime           int                `bson:"start_time,omitempty"`
		EndTime             int                `bson:"end_time,omitempty"`
		StartKm             int                `bson:"start_km,omitempty"`
		EndKm               int                `bson:"end_km,omitempty"`
		TotalDistance       int                `bson:"total_distance,omitempty"`
		VehicleType         string             `bson:"vehicle_type,omitempty"`
		BookingType         string             `bson:"booking_type,omitempty"`
		CreatedTime         time.Time          `bson:"created_time,omitempty"`
		CarbonEmissionSaved int                `bson:"carbon_emission_saved,omitempty"`
		StartingStation     struct {
			ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
			Name        string             `bson:"name,omitempty"`
			Description string             `bson:"description,omitempty"`
			Location    struct {
				Type        string    `bson:"type,omitempty"`
				Coordinates []float64 `bson:"coordinates,omitempty"`
			} `bson:"location,omitempty"`
		} `bson:"starting_station,omitempty"`
		EndingStation struct {
			ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
			Name     string             `bson:"name,omitempty"`
			Location struct {
				Type        string    `bson:"type,omitempty"`
				Coordinates []float64 `bson:"coordinates,omitempty"`
			} `bson:"location,omitempty"`
		} `bson:"ending_station,omitempty"`
		CouponCode        string    `bson:"coupon_code,omitempty"`
		Discount          int       `bson:"discount,omitempty"`
		GreenPoints       int       `bson:"green_points,omitempty"`
		CarbonSaved       int       `bson:"carbon_saved,omitempty"`
		City              string    `bson:"city,omitempty"`
		RideTimeRemaining int       `bson:"ride_time_remaining,omitempty"`
		UpdateTime        time.Time `bson:"update_time,omitempty"`
		TimeBooked        int       `bson:"time_booked,omitempty"`
		UpdatedAt         time.Time `bson:"updated_at,omitempty"`
	} `json:"booking,omitempty"`
	Plan struct {
		ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		Name        string             `bson:"name,omitempty"`
		City        string             `bson:"city,omitempty"`
		VehicleType string             `bson:"vehicle_type,omitempty"`
		Type        string             `bson:"type,omitempty"`
		Price       int                `bson:"price,omitempty"`
	} `bson:"plan,omitempty"`
	PlanType string `bson:"planType,omitempty"`
	Price    int    `bson:"price,omitempty"`
	Profile  struct {
		ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		Name string             `bson:"name,omitempty"`
	} `bson:"profile,omitempty"`
}

func GetVehicleData(id int) ([]VehicleData, error) {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "deviceId", Value: 310823018}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "booking"},
					{Key: "localField", Value: "deviceId"},
					{Key: "foreignField", Value: "device_id"},
					{Key: "as", Value: "booking"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$booking"}}}},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "planId", Value: bson.D{{Key: "$toString", Value: "$booking.plan._id"}}},
					{Key: "userID", Value: bson.D{{Key: "$toObjectId", Value: "$booking.profile_id"}}},
					{Key: "plan", Value: "$booking.plan"},
					{Key: "planType", Value: "$booking.booking_type"},
					{Key: "price", Value: "$booking.price"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "users"},
					{Key: "localField", Value: "userID"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "profile"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}},
	}
	cursor, err := trestCommon.Aggregate(pipeline, "iotBike")
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var result []VehicleData
	for cursor.Next(context.TODO()) {
		var profile VehicleData
		if err = cursor.Decode(&profile); err != nil {
			trestCommon.ECLog1(err)
			continue
		}
		result = append(result, profile)
	}
	return result, nil
}

func ImmobilizeDevice(deviceId int) error {
	filter := bson.M{"device_id": deviceId}
	repo := bikeDevice.NewRepository("bikeDevice")
	device, err := repo.FindOne(filter, bson.M{})

	if err != nil {
		return err
	}
	repoIot := iotbike.NewProfileRepository("iotBike")
	iotDev, _ := repoIot.Find(bson.M{"deviceId": deviceId}, bson.M{})
	if device.DeviceID == 0 {
		return errors.New("device not found")
	}
	set := bson.M{"immobilized": !device.Immobilized}

	if len(iotDev) == 1 && iotDev[0].DeviceId != 0 && iotDev[0].Type == "moto" {
		if device.Immobilized {
			motog.ImmoblizeDevice(0, iotDev[0].Name)
		} else {
			motog.ImmoblizeDevice(1, iotDev[0].Name)

		}
	}
	_, err = repo.UpdateOne(filter, bson.M{"$set": set})
	return err
}
