app = angular.module('UnionCarWebsiteApp', [])

app.config ['$routeProvider', ($routeProvider) ->
	$routeProvider.when '/',
		templateUrl: 'views/main.html'
		controller: 'MainCtrl'
	$routeProvider.otherwise
		redirectTo: '/']
