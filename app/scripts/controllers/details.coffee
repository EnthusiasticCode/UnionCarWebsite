'use strict'

angular.module('UnionCarWebsiteApp')
	.controller 'DetailsCtrl', ($scope, $routeParams, carApi) ->
		$scope.id = $routeParams.id
		carApi.get $scope.id, (car) ->
			$scope.car = car
