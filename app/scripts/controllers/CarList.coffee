'use strict'

angular.module('UnionCarWebsiteApp')
.controller 'CarListCtrl', ['$scope', ($scope) ->
	$scope.cars = [
		{
			id: 1,
			brand: 'Ferrari',
			model: 'Testarossa',
			km: 110000,
			date: new Date(2003, 3),
			price: 230000,
			images: [
				'http://placehold.it/300&text=Ferrari'
			]
			},
			{
				id: 1,
				brand: 'Ferrari',
				model: 'Testarossa',
				km: 110000,
				date: new Date(2003, 3),
				price: 20000,
				images: [
					'http://placehold.it/300&text=Ferrari'
				]
			}
		]
	]
