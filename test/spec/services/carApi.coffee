'use strict'

describe 'Service: car', () ->

  # load the service's module
  beforeEach module 'UnionCarWebsiteApp'

  # instantiate service
  car = {}
  beforeEach inject (_car_) ->
    car = _car_

  it 'should do something', () ->
    expect(!!car).toBe true;
