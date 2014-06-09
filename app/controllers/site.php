<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Site extends CI_Controller {

	public function index($car_id=-1)
	{
		$cars = null;
		if ( $this->agent->is_robot() ) {
			$this->load->model('Cars');
			if ( $car_id >= 0 ) {
				$cars = array( $this->Cars->get($car_id) );
				if ( count($cars) == 0 )
					$cars[] = new Cars();
			} else {
				$cars = $this->Cars->all();
			}
		}
		$this->load->view('site', array( 'cars' => $cars ));
	}

	public function views($view_name)
	{
		$this->load->model('Cars');
		$this->load->view($view_name, array( 'cars' => array( new Cars() ) ));
	}

	public function mail() {
		$this->load->library('email');

		$this->email->from($this->input->post('sender'));
		$this->email->to('info@unioncar.it');

		$this->email->subject('[site] contact form');
		$this->email->message($this->input->post('text'));

		if ($this->email->send()) {
			$this->output
				->set_content_type('application/json')
				->set_output(json_encode(array('status' => 'ok')));
		} else {
			$this->output
				->set_header("HTTP/1.1 500 Internal Server Error")
				->set_content_type('application/json')
				->set_output(json_encode(array('error' => $this->email->print_debugger())));
		}
	}

}
