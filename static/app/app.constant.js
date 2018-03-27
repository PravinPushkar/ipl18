'use strict';

var app = angular.module('ipl');

/**
 * Constants
 * 
 * Contains all the constants for the angular application.
 */
app.constant('INumberPattern', /^[i|I][0-9]{6}$/)
    .constant('aliasPattern', /^\w{0,20}$/);