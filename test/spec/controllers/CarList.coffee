'use strict'

describe 'Controller: CarListCtrl', () ->

  # load the controller's module
  beforeEach module 'UnionCarWebsiteApp'

  CarListCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    CarListCtrl = $controller 'CarListCtrl', {
      $scope: scope
    }

  it 'should attach a list of cars to the scope', () ->
    expect(scope.cars.length).toBe 1;
