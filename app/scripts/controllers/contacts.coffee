'use strict'

window.loadContactsMaps = ->
	map = new google.maps.Map(document.getElementById("map-canvas"), {
			zoom: 8,
			center: new google.maps.LatLng(-34.397, 150.644),
			mapTypeId: google.maps.MapTypeId.ROADMAP
		})

angular.module('UnionCarWebsiteApp')
	.controller 'ContactsCtrl', ($scope, $timeout) ->
		# Loading Google Maps API if needed
		# unless google?.maps?
		# 	script = document.createElement("script")
		# 	script.type = "text/javascript"
		# 	script.src = "http://maps.googleapis.com/maps/api/js?client=gme-yourclientid&sensor=false&callback=loadContactsMaps"
		# 	document.body.appendChild(script);
		# else
		# 	loadContactsMaps()

		$scope.mail =
			sender: null
			text: null
			send: ->
				$scope.mail.status = "sending"
				# TODO substitute with API ajax call
				$timeout (->
					console.log $scope.mail.sender, $scope.mail.text
					$scope.mail.status = "sent"), 2000
			status: "form"
