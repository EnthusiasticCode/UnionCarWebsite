<?php if (!defined('BASEPATH')) exit('No direct script access allowed');

class MY_Model extends CI_Model
{
	var $id = 0;

	var $_model_name = '';
	var $_db_table = '';
	var $_split_char = ':';
	var $_target = Null;
	var $_aliases = array();
	var $_rules = array();

	var $_generic_filter_key = 'search';
	var $_generic_filter_fields = array();

	var $_render_string_defaul = 'n/d';
	var $_render_integer_defaul = 'n/d';
	var $_render_bool_true = 'Si';
	var $_render_bool_false = 'No';

	var $_render_options_open = '<div class="options" style="display:none">';
	var $_render_options_close = '</div>';

	function __construct()
	{
		parent::__construct();
		$this->_model_name = get_class($this);
		$this->_db_table = strtolower($this->_model_name);
	}

	function populate($values, $set = false)
	{
		$any_set = !$set;
		if (empty($values)) {
			$values = &$this;
			if ($set) $set = 'all';
		}
		foreach ($values as $key => $val) {
			// Adjust key
			if (strpos($key, $this->_db_table.'_') === 0)
				$key = substr($key, strlen($this->_db_table.'_'));
			// Check if key is valid
			if (property_exists($this->_model_name, $key)
				AND strpos($key, '_') !== 0)
			{
				$t = gettype($this->{$key});
				if ($t == 'boolean')
					$val = $val === 'false' ? false : (bool)$val;
				else
					settype($val, $t);
				// Set database value
				if ($set AND $key != 'id' AND ($set === 'all' OR $this->{$key} !== $val)) {
					$this->db->set($key, $val);
					$any_set = true;
				}
				// Reflect changes
				$this->{$key} = $val;
			}
		}
		return $any_set;
	}

	function get($id)
	{
		$this->db->select('*');
		$this->db->where('id', $id);
		$res = $this->db->get($this->_db_table);
		if ($res->num_rows() > 0)
		{
			$this->populate($res->row());
		}
		else
		{
			$this->id = 0;
			// TODO proper reset
			// Maybe saving initial object state and restoring that
		}
		return $this;
	}

	function all()
	{
		return $this->db->get($this->_db_table)->result();
	}

	function random($count = 1)
	{
		return $this->db->order_by('id', 'random')->limit($count)->get($this->_db_table)->result();
	}

	function count_all()
	{
		return $this->db->count_all($this->_db_table);
	}

	function save($values = Null, $with_id = Null)
	{
		// Update form array
		$values = $this->_save_internal($values);
		if ($with_id)
		{
			$this->db->set('id', $with_id);
		}
		if ($this->populate($values, 'all') && $this->db->insert($this->_db_table))
		{
			$this->id = $this->db->insert_id();
			return $this->id;
		}
		return !empty($values);
	}

	function _save_internal($values) { return $values; }

	function update($values = Null)
	{
		// Update form array
		$values = $this->_update_internal($values);
		if ($this->populate($values, true)) {
			$this->db->where('id', $this->id);
			return $this->db->update($this->_db_table);
		}
		return !empty($values);
	}

	function _update_internal($values) { return $values; }

	function save_or_update($values, $where)
	{
		$this->id = 0;
		$this->filter($where, 1);
		return $this->id > 0
			? $this->update($values)
			: $this->save($values);
	}

	function delete($id = 0)
	{
		if ($id == 0)
			$id = $this->id;
		if ($id > 0)
		{
			$this->db->delete($this->_db_table, array('id' => $id));
		}
	}

	function filter($filters_values = array(), $pagesize = 0) {

		// Pagesize as first parameter
		if (is_numeric($filters_values))
		{
			$pagesize = $filters_values;
			$filters_values = NULL;
		}

		// Se non sono specificati valori uso il segmento
		if (empty($filters_values))
			$filters_values = $this->uri->uri_to_assoc();

		// Do priorita' alle variabili POST
		if (!empty($_POST))
			$filters_values = array_merge($filters_values, $_POST);

		// Remove id as filter
		if (isset($filters_values['id']))
			unset($filters_values['id']);

		// Inizializzazione query
		$this->db->start_cache();
		$this->_filter_from($filters_values, $pagesize);

		// Impostazione filtri
		$cur_page = 0;
		$sort = array();
		if (!empty($filters_values))
		{
			// Getting page
			if (isset($filters_values['page']))
			{
				if (is_numeric($filters_values['page']))
					$cur_page = (int) $filters_values['page'];
			}

			// Sorting
			if(isset($filters_values['sortby']))
			{
				$sortby = explode($this->_split_char, $filters_values['sortby']);
				$sortorder = isset($filters_values['sortorder'])
					? explode($this->_split_char, $filters_values['sortorder'])
					: array('desc');
				$dup = count($sortby) - count($sortorder);
				while ($dup-- > 0) $sortorder[] = $sortorder[count($sortorder) - 1];
				for ($i = 0; $i < count($sortby); $i++) if (property_exists($this->_model_name, $sortby[$i])) {
					$this->db->order_by($sortby[$i], $sortorder[$i]);
					$sort[$sortby[$i]] = $sortorder[$i];
				}
			}

			// Generic filtering
			if (!empty($this->_generic_filter_fields)
			AND isset($filters_values[$this->_generic_filter_key])) {
				$val = $filters_values[$this->_generic_filter_key];
				$count = 0;
				$query = '(';
				foreach ($this->_generic_filter_fields as $column => $func) {
					$query .= ($count ? ' OR ' : '');
					$wildchar = ($func == 'LIKE' ? '%' : '');
					$query .= $column . ' ' . $func . ' ' . (is_string($this->{$column})
						? $this->db->escape($wildchar.$val.$wildchar)
						: $val);
					$count++;
				}
				$query .= ')';
				$this->db->where($query);
			}

			// Automatic filtering
			foreach ($filters_values as $key => $value)
			{
				// Key aliasing
				$column = $key;
				if (isset($this->_aliases[$key]))
					$column = $this->_aliases[$key];

				// Filter translation
				if (strpos($column, '_') !== 0 AND property_exists($this->_model_name, $column) AND !empty($value))
				{
					// Getting filter value
					if (!is_array($value))
					{
						$value = is_string($value) ? explode($this->_split_char, urldecode($value)) : array($value);
					}
					// Query generation
					$query = '(';
					$count = 0;
					$concat = property_exists($this->_model_name,'_'.$key.'_concat') ? $this->{'_'.$key.'_concat'} : 'OR';
					$func = property_exists($this->_model_name,'_'.$key.'_func')
						? $this->{'_'.$key.'_func'}
						: (is_string($this->{$column}) ? 'LIKE' : '=');
					$wildchar = $func == 'LIKE' ? '%' : '';
					foreach ($value as $val)
					{
						if ($count > 0) $query .= ' ' . $concat . ' ';
						$query .= $column . ' ' . $func . ' ' . (is_string($this->{$column})
							? $this->db->escape($wildchar.$val.$wildchar)
							: $val);
						$count++;
					}
					$query .= ')';
					// Add to query
					$this->db->where($query);
				}
				else unset($filters_values[$key]);
			}
		}

		// Post filtering
		$callback = '_'.$this->_db_table.'_filter';
		if($this->_target AND method_exists($this->_target, $callback))
			$this->_target->$callback();

		// Stop cache to count results before limit
		$this->db->stop_cache();

		// Pre-pageing count
		$callback = '_'.$this->_db_table.'_filter_count';
		if($this->_target AND method_exists($this->_target, $callback))
			$this->_target->$callback($filters_values, $this->db->count_all_results(), $pagesize, $sort);

		// Setup pagination
		if ($pagesize > 0)
		{
			$this->db->limit($pagesize, $cur_page);
		}

		// Call target refine method
		$callback = '_'.$this->_db_table.'_filter_refine';
		if($this->_target AND method_exists($this->_target, $callback))
			$this->_target->$callback();

		// Run query
		$res = $this->db->get();
		$this->db->flush_cache();

		// Populate if only 1 record requested
		if ($pagesize == 1) {
			if ($res->num_rows() > 0) $this->populate($res->row());
			else $this->id = 0; // TODO reset all to default
			return $this;
		}

		return $res->result();
	}

	function _filter_from($filters_values, $pagesize)
	{
		$this->db->from($this->_db_table);
	}

	function validate($rules = array())
	{
		$this->load->library('form_validation');
		if (empty($rules))
			$rules = &$this->_rules;
		// Generate validation configuration
		foreach ($this as $key => $value)
		if (strpos($key, '_') !== 0
			AND !(isset($rules[$key]) AND $rules[$key] == 'hide'))
		{
			// Add validation rule
			$this->form_validation->set_rules(
				$this->_db_table.'_'.$key,
				$key,
				(property_exists($this->_model_name,'_'.$key.'_rule')
					? $this->{'_'.$key.'_rule'}
					: (isset($rules[$key]) ? $rules[$key] : ''))
			);
		}
		// Return validation result
		return $this->form_validation->run();
	}

	function render($field, $edit = true, $params = array())
	{
		// Default parameters
		if (!is_array($params))
			$params = array('type' => $params);
		else if (!isset($params['type']))
			$params['type'] = 'text';
		//
		if (property_exists($this->_model_name, $field)
			AND strpos($field, '_') !== 0
			AND !(isset($this->_rules[$field]) AND $this->_rules[$field] == 'hide'))
		{
			$name = $this->_db_table.'_'.$field;
			// Autocomplete
			$autocomplete = false;
			if ($edit AND method_exists($this->_model_name, '_'.$field.'_autocomplete'))
				$autocomplete = $this->{'_'.$field.'_autocomplete'}();
			// Render values field
			if (property_exists($this->_model_name, '_'.$field.'_values')) {
				if ($edit) {
					$render =
						'<select name="'.$name
						.'" title="'.$field
						.'" id="'.$name.'"'
						.$this->_params_string($params).'>';
					foreach($this->{'_'.$field.'_values'} as $v => $n) $render .=
							'<option value="'.$v.'"'.
							($this->input->post($name) === false
								? ($this->{$field} == $v ? ' selected="selected"' : '')
								: ($this->input->post($name) == $v ? ' selected="selected"' : ''))
							.'>'
							.$n
							.'</option>';
					$render .= '</select>';
				} else $render = $this->{'_'.$field.'_values'}[$this->{$field}];
			}
			// Render by field type
			else switch (gettype($this->{$field}))
			{
			case 'boolean':
				$render = $edit
					? '<input type="hidden" name="'.$name.'" value="false" />'
						.'<input type="checkbox" name="'.$name
						.'" id="'.$name.'"'
						.($this->input->post($name) === false
							? ($this->{$field} ? ' checked="checked"' : '')
							: ($this->input->post($name) == 'true' ? ' checked="checked"' : ''))
						.$this->_params_string($params)
						.' title="'.$field
						.'" value="true" />'
					: ($this->{$field} ? $this->_render_bool_true : $this->_render_bool_false);
				break;
			case 'integer':
				if (!$edit) { $render = ($this->{$field} ? $this->{$field} : $this->_render_integer_defaul); break; }
				else isset($params['class']) ? $params['class'] .= ' number' : $params['class'] = 'numeric';
			case 'string':
				if ($params['type'] == 'textarea') {
				$render = $edit
					? '<textarea name="'.$name
						.'" title="'.$field
						.'" id="'.$name.'"'
						.$this->_params_string($params).'>'
						.($this->input->post($name) === false ? $this->{$field} : $this->input->post($name))
						.'</textarea>'
					: ($this->{$field} ? $this->{$field} : $this->_render_string_defaul); break; }
			default:
				$render = $edit
					? '<input type="'.$params['type']
						.'" name="'.$name
						.'" id="'.$name
						.'" value="'.($this->input->post($name) === false ? $this->{$field} : $this->input->post($name))
						.'" title="'.$field
						.'"'.$this->_params_string($params).' />'
					: ($this->{$field} ? $this->{$field} : $this->_render_string_defaul);
			}
			// Rendering
			return $render.(($autocomplete OR $this->input->post($name) !== false) ? (
				$this->_render_options_open
				.($autocomplete ? $this->_render_autocomplete($autocomplete) : '')
				.($this->input->post($name) !== false ? $this->_render_default($this->{$field}) : '')
				.(form_error($name) ? $this->_render_error(form_error($name)) : '')
				.$this->_render_options_close
				) : '');
		}
		return '';
	}

	function _render_autocomplete($data) {
		return '<div class="autocomplete">'.implode(',',$data).'</div>';
	}

	function _render_default($data) {
		return '<div class="default">'.$data.'</div>';
	}

	function _render_error($data) {
		return '<div class="error">'.$data.'</div>';
	}

	function _params_string($params)
	{
		unset($params['type']);
		$res = '';
		foreach ($params as $name => $val)
			$res .= ' '.$name.'="'.$val.'"';
		return $res;
	}

	function serialize($map = NULL, $values = NULL)
	{
		// Serialize this object
		if (!isset($values))
			$values =& $this;
		// Serialization result
		$serialization = '';
		// Map interpretation
		$map_usage = false;
		if (is_array($map))
			$map_usage = isset($map['usage']) ? $map['usage'] : 'specific';
		// Serialization loop
		foreach ($values as $key => $value)
		{
			// Check if key is valid
			if (property_exists($this->_model_name, $key)
				AND strpos($key, '_') !== 0)
			{
				// Check if name needs to be mapped
				if ($map_usage)
				{
					if ($map_usage == 'include' AND !in_array($key, $map)) continue;
					if ($map_usage == 'exclude' AND  in_array($key, $map)) continue;
					if (isset($map[$key])) $key = $map[$key];
					if ($key === false) continue;
				}
				// Serialize key-value
				$serialization .= $this->_serialize_node($key, $this->_serialize_value($value));
			}
		}
		// Ending serialization
		if (!empty($serialization))
		{
			$root = $this->_db_table;
			if (is_array($map) AND isset($map[$root])) $root = $map[$root];
			$serialization = $this->_serialize_node($root, "\n".$serialization);
		}
		// Returning
		return $serialization;
	}

	// Accepted input:
	// array(
	// 	'tag' => 'casa',
	// 	'@attr' => '$field'
	// )
	function _serialize_node($node, $value = NULL, $context = NULL)
	{
		$tag = $this->_db_table;
		$attributes = '';
		if ($context == NULL)
			$context =& $this;
		// Node with attributes
		if (is_array($node)) foreach ($node as $key => $val)
		{
			if ($key == 'tag')
				$tag = $val;
			if ($key == 'value')
				$value = $val;
			elseif (strpos($key, '@') === 0)
			{
				$attributes .= ' '
				.strip_tags(substr($key, 1))
				.'="'
				.(strpos($val, '$') === 0
					? $context->{substr($val, 1)}
					: $val)
				.'"';
			}
		}
		else $tag = $node;
		// Simple node
		$tag = strip_tags($tag);
		return $value === NULL ? "<".$tag.$attributes."/>\n" : '<'.$tag.$attributes.">".$value."</".$tag.">\n";
	}

	function _serialize_value($value, $strip = true)
	{
		if (is_string($value))
			return '<![CDATA['.strip_tags($value).']]>';
		return $value;
	}
}
?>