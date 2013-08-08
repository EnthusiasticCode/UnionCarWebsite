app = angular.module('UnionCarWebsiteApp', ['ngRoute', 'ngResource', 'angular-flexslider'])

app.config ($routeProvider, $locationProvider, conf) ->
	$routeProvider.when '/',
		templateUrl: "#{conf.siteUrl}/views/carlist",
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: "#{conf.siteUrl}/views/details",
		controller: 'DetailsCtrl'
	.when '/contacts',
		templateUrl: "#{conf.siteUrl}/views/contacts",
		controller: 'ContactsCtrl'
	.when '/contacts/:carId',
		templateUrl: "#{conf.siteUrl}/views/contacts",
		controller: 'ContactsCtrl'
	.otherwise
		redirectTo: '/'

	$locationProvider.html5Mode(yes).hashPrefix('!')
