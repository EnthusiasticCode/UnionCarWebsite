app = angular.module('UnionCarWebsiteApp', ['ngRoute', 'ngResource', 'angular-flexslider'])

app.config ($routeProvider, $locationProvider, conf) ->
	$routeProvider.when '/',
		templateUrl: "#{conf.baseUrl}/views/carlist",
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: "#{conf.baseUrl}/views/details",
		controller: 'DetailsCtrl'
	.when '/contacts',
		templateUrl: "#{conf.baseUrl}/views/contacts",
		controller: 'ContactsCtrl'
	.when '/contacts/:carId',
		templateUrl: "#{conf.baseUrl}/views/contacts",
		controller: 'ContactsCtrl'
	.otherwise
		redirectTo: '/'

	$locationProvider.html5Mode(yes).hashPrefix('!')
