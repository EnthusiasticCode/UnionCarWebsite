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

parseCar = (car) ->
	car.brand = car.brand.toLowerCase()

	car.images = [
		car.image_url_1 or 'nopicture.jpg',
		car.image_url_2,
		car.image_url_3,
		car.image_url_4
	]
	delete car.image_url_1
	delete car.image_url_2
	delete car.image_url_3
	delete car.image_url_4

	car.price = parseInt(car.price)

	if car.date == '0000-00-00'
		car.date = null
	else
		car.date = car.date.split('-')
		car.date = new Date(car.date[0], car.date[1] - 1, car.date[2])

	return car

# Adding car provider
angular.module('UnionCarWebsiteApp')
	.factory 'carApi', carFactory
