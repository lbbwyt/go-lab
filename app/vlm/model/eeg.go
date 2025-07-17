package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MultiModelData struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	TimeStep     string             `bson:"timestep"`
	DeviceInfoId string             `bson:"device_info_id"`
	EegChannel1  float64            `bson:"eeg_channel1"`
	EegChannel2  float64            `bson:"eeg_channel2"`

	EegChannel3 float64 `bson:"eeg_channel3"`
	EegChannel4 float64 `bson:"eeg_channel4"`

	EogChannel1    float64 `bson:"eog_channel1"`
	EogChannel2    float64 `bson:"eog_channel2"`
	TriggerChannel int     `bson:"trigger_channel"`
	PpgRed         int     `bson:"ppg_red"`

	PpgInfrared int `bson:"ppg_infrared"`

	ImuAccX float64 `bson:"imu_acc_x"`
	ImuAccY float64 `bson:"imu_acc_y"`
	ImuAccZ float64 `bson:"imu_acc_z"`

	ImuGyroX float64 `bson:"imu_gyro_x"`
	ImuGyroY float64 `bson:"imu_gyro_y"`
	ImuGyroZ float64 `bson:"imu_gyro_z"`

	ImagePath string `bson:"image_path"`
	Label     string `bson:"label"`
}
