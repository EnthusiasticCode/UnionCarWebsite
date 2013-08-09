<?php if (!defined('BASEPATH')) exit('No direct script access allowed');

class Cars extends MY_Model {

	var $fuel_type = '&ndash;'; // char(1) NOT NULL,
	var $pubblic = true; // tinyint(1) NOT NULL,
	var $gear_type = '&ndash;'; // varchar(35) NOT NULL,
	var $car_type = '&ndash;'; // varchar(35) NOT NULL,
	var $km = 0; // int(11) NOT NULL,
	var $engine_size = 0; // int(11) NOT NULL,
	var $internal_code = ''; // varchar(20) NOT NULL,
	var $infocar = ''; // varchar(12) NOT NULL,
	var $external_color = '&ndash;'; // varchar(15) NOT NULL,
	var $registration_date = ''; // date NOT NULL,
	var $brand = '&ndash;'; // varchar(20) NOT NULL,
	var $model = '&ndash;'; // varchar(40) NOT NULL,
	var $model_description = ''; // varchar(100) NOT NULL,
	var $power_horses = 0; // int(11) NOT NULL,
	var $power_kw = 0; // int(11) NOT NULL,
	var $price = 0; // int(11) NOT NULL,
	var $image_url_1 = ''; // varchar(15) NOT NULL,
	var $image_url_2 = ''; // varchar(15) NOT NULL,
	var $image_url_3 = ''; // varchar(15) NOT NULL,
	var $image_url_4 = ''; // varchar(15) NOT NULL,
	var $image_url_5 = ''; // varchar(15) NOT NULL,
	var $image_url_6 = ''; // varchar(15) NOT NULL,
	var $image_url_7 = ''; // varchar(15) NOT NULL,
	var $image_url_8 = ''; // varchar(15) NOT NULL,

}
