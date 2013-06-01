app = angular.module('UnionCarWebsiteApp', [])

app.config ['$routeProvider', ($routeProvider) ->
	$routeProvider.when '/',
		templateUrl: 'views/carlist.html'
		controller: 'CarListCtrl'
	.when '/undefined',
		templateUrl: 'views/undefined.html',
		controller: 'UndefinedCtrl'
	.when '/details/:id',
		templateUrl: 'views/details.html',
		controller: 'DetailsCtrl'
	.otherwise
		redirectTo: '/']
