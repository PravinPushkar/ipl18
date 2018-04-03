'use strict';

var app = angular.module('ipl');

/**
 * ImageSelect Directive
 * 
 * Directive to select the image uploading/uploaded.
 */
app.directive('imageSelect', ['$timeout', 'imageReader', function ($timeout, imageReader) {
    return {
        scope: {
            ngModel: '='
        },
        link: function ($scope, element) {
            function getFile(file) {
                imageReader.readAsDataUrl(file, $scope)
                    .then(function (result) {
                        $timeout(function () {
                            $scope.ngModel = result;
                        });
                    });
            }
            element.on('change', function (event) {
                var file = (event.srcElement || event.target).files[0];
                getFile(file);
            });
        }
    };
}]);