'use strict';

var app = angular.module('ipl');

/**
 * Constants
 * 
 * Contains all the constants for the angular application.
 */
app.constant('INumberPattern', /^[0-9]{6}$/)
    .constant('aliasPattern', /^\w{0,20}$/);