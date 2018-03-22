'use strict';

var app = angular.module('ipl');

/**
 * Edit Profile Controller
 * 
 * Controller for the edit profile page.
 */
app.controller('editProfileController', function ($http, $mdToast, $scope, urlService, utilsService) {
    var vm = this;

    vm.edit = edit;

    // Get the file name of image uploaded
    document.querySelector('input[id="profilePic"]').onchange = function () {
        vm.imageSelectedName = document.getElementById('profilePic').files[0].name;
    };

    // $scope.$on('fileProgress', function (e, progress) {
    //     vm.progress = progress.loaded / progress.total;
    // });

    // Send editable data to back-end
    function edit() {
        var data = {};
        if (vm.password !== vm.confirmPassword) {
            utilsService.showToast({
                text: 'Password and Confirm Password do not match',
                hideDelay: 3000,
                isError: true
            });
            return;
        }
        if (vm.password !== '' || vm.password !== undefined || vm.password !== null) {
            data.password = vm.password;
        }
        if (vm.alias !== '' || vm.alias !== undefined || vm.alias !== null) {
            data.alias = vm.alias;
        }
        if (vm.profilePic !== '' || vm.profilePic !== undefined || vm.profilePic !== null) {
            data.profilePic = vm.profilePic;
        }
        $http.put(urlService.userProfile)
            .then(function (res) {
                // vm.userData = {
                //     firstName: res.body.firstName
                // };
                console.log('success');
            }, function () {
                console.log('error');
            });
    }
});