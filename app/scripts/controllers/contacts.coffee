'use strict'

angular.module('UnionCarWebsiteApp').controller 'ContactsCtrl', ($scope, $routeParams, carApi) ->
	$scope.mail =
		status: "form"
		data:
			sender: null
			text: null
		send: ->
			$scope.mail.status = "sending"
			$.ajax
				type: 'POST'
				data: $scope.mail.data
				url: 'http://www.unioncar.it/index.php/site/mail'
				success: -> $scope.$apply ->
					$scope.mail.status = "sent"
				error: -> $scope.$apply ->
					$scope.mail.status = "error"

	# Load default message if needed
	if $routeParams.carId
		carApi.get $routeParams.carId, (car) ->
			$scope.mail.data.text = """Gentile Union Car,

				Vorrei maggiori informazioni in merito alla '#{car.brand} #{car.model}'#{(car.date&&(' del '+car.date.getFullYear()))||''}.
				In particolare:\n"""
