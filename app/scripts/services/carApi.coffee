'use strict';

carFactory = ($resource) ->
	# Service logic
	Car = $resource('/cgi-bin/api/car/:carId', {carId: '@id'})
	allCars = null
	allGets = {}

	# Public API here
	{
		resource: Car
		all: (callback) ->
			if allCars?
				callback(allCars, null)
			else
				Car.query (data, response) ->
					allCars = (parseCar(d) for d in data)
					callback(data, response)
			allCars
		get: (carId, callback) ->
			car = allGets[carId]
			if car?
				callback(car, null)
			else
				Car.get {carId: carId}, (data, response) ->
					allGets[carId] = car = parseCar(data)
					callback(data, response)
			car
	}
carFactory.$inject = ['$resource']

fuelTypes =
	'B': 'benzina'
	'D': 'disel'
	'G': 'gpl'
	'M': 'metano'
	'I': 'idrogeno'

parseCar = (car) ->
	car.brand = car.brand.toLowerCase()
	car.model = car.model.replace /\s+$/, ''

	car.images = []
	for i in [1..8] when c = car["image_url_#{i}"]
		car.images.push c
		delete car["image_url_#{i}"]

	car.pubblic = car.pubblic == '1' if car.pubblic?
	car.engine_size = parseInt(car.engine_size) if car.engine_size?
	car.price = parseInt(car.price) if car.price?
	car.km = parseInt(car.km) if car.km?
	car.power_kw = parseInt(car.power_kw) if car.power_kw?
	car.power_horses = parseInt(car.power_horses) if car.power_horses?
	car.fuel_type = fuelTypes[car.fuel_type] if car.fuel_type?

	if car.registration_date == '0000-00-00'
		car.registration_date = null
	else
		car.registration_date = new Date(car.registration_date)

	return car

# Adding car provider
angular.module('UnionCarWebsiteApp')
	.factory 'carApi', carFactory
