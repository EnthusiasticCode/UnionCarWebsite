'use strict'

window.loadContactsMaps = ->
	address = $('.vcard.address').text().replace(/\n\s*/g, ' ')
	map = new google.maps.Map(document.getElementById("map-canvas"), {
			zoom: 15,
			center: new google.maps.LatLng(45.494296, 12.091114),
			mapTypeId: google.maps.MapTypeId.ROADMAP
		})
	geocoder = new google.maps.Geocoder()
	geocoder.geocode { 'address': address }, (results, status) ->
		if status == google.maps.GeocoderStatus.OK
			map.setCenter(results[0].geometry.location)
			marker = new google.maps.Marker
				map: map
				position: results[0].geometry.location

angular.module('UnionCarWebsiteApp')
	.controller 'ContactsCtrl', ($scope, $http) ->
		# Loading Google Maps API if needed
		unless google?.maps?
			script = document.createElement("script")
			script.type = "text/javascript"
			script.src = "https://maps.googleapis.com/maps/api/js?key=AIzaSyAmeH3LafBcub934HRIM8jNtg_xTNZFQ5Y&sensor=false&callback=loadContactsMaps"
			document.body.appendChild(script)
		else
			loadContactsMaps()

		$scope.mail =
			status: "form"
			data:
				sender: null
				text: null
			send: ->
				$scope.mail.status = "sending"
				$http.post('/cgi-bin/api/mail', $scope.mail.data).success(-> $scope.mail.status = "sent")
