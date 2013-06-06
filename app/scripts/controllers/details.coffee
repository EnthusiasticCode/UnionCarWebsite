'use strict'

angular.module('UnionCarWebsiteApp')
	.controller 'DetailsCtrl', ($scope, $routeParams) ->
		$scope.id = $routeParams.id
		$scope.car = {
			id: 1,
			brand: 'Ferrari',
			model: 'Testarossa',
			km: 110000,
			date: new Date(2003, 3),
			price: 230000,
			images: [
				'http://placehold.it/300&text=Ferrari'
			]
		}

		$('.flexslider').flexslider({
			animation: "slide",
			controlNav: "thumbnails"
		})
