'use strict'

angular.module('UnionCarWebsiteApp')
	.controller 'CommerciantiCtrl', ($scope, $location) ->
		$scope.login =
			email: null
			password: null
			login: ->
				$.ajax
					type: 'POST'
					url: '/cgi-bin/api/login'
					data:
						email: $scope.login.email
						password: $scope.login.password
					success: (data) -> $scope.$apply ->
						$location.path('/') if data?.login is "ok"
						$scope.login.error = yes
