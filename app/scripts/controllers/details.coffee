'use strict'

angular.module('UnionCarWebsiteApp')
  .controller 'DetailsCtrl', ($scope, $routeParams) ->
    $scope.id = $routeParams.id
