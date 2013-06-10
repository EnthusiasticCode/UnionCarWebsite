<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Api extends CI_Controller {

	public function car($id = null) {
		// $this->load->library('unzip');
		// $this->load->library('csvreader');

		$this->output
			->set_content_type('application/json')
			->set_output(json_encode(array('foo' => 'bar', 'id' => $id)));
	}

}