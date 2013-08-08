<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Site extends CI_Controller {

	public function index()
	{
		$this->load->view('site');
	}

	public function views($view_name)
	{
		$this->load->view($view_name);
	}

}
