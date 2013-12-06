<!doctype html>
<!--[if lt IE 9]>      <html class="no-js lt-ie9"> <![endif]-->
<!--[if gt IE 8]><!--> <html class="no-js"> <!--<![endif]-->
  <head>
    <title>Union Car - Auto Nuove e Usate</title>
    <meta name="description" content="Auto nuove e usate a Mirano Venezia Italia.">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- Place favicon.ico and apple-touch-icon.png in the root directory -->
    <link rel='stylesheet' href='http://fonts.googleapis.com/css?family=Lato:100,300,700' type='text/css'>
    <link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/select2/3.4.0/select2.min.css">
    <link rel="stylesheet" href="<?php echo base_url(APPPATH.'/styles/main.css'); ?>">
    <!--[if lt IE 9]>
    <style type="text/css" src="/styles/ie.css"></style>
    <script src="//cdnjs.cloudflare.com/ajax/libs/html5shiv/3.6.2/html5shiv.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/json3/3.2.4/json3.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/respond.js/1.1.0/respond.min.js"></script>
    <![endif]-->
  </head>
  <body ng-app="UnionCarWebsiteApp">
    <!--[if lt IE 9]>
      <p class="chromeframe">Stai usando un browser obsoleto. <a target="_blank" href="http://browsehappy.com/">Aggiorna il tuo browser oggi</a> o <a target="_blank" href="http://www.google.com/chromeframe/?redirect=true">install Google Chrome Frame</a> per un'esperienza migliore.</p>
    <![endif]-->

    <div id="site" class="row">
      <div id="page" class="small-12 columns">

        <header id="page-header" class="row">
          <ul class="inline-list brand-list hide-for-small">
            <li class="header-brand-opel"></li>
            <li class="header-brand-mazda"></li>
            <li class="header-brand-subaru"></li>
            <li class="header-brand-seat"></li>
          </ul>
          <div class="large-4 small-12 columns">
            <h1 id="page-header-logo"><a class="link" href="/">Union Car</a></h1>
            <ul class="vcard logo-vcard address contacts hide-for-small">
              <li>Tel: <span class="phone-number">+39 041 570 3547</span></li>
            </ul>
          </div>
          <div class="large-5 small-12 columns hide-for-small">
            <h2 id="page-header-slogan" class="subheader">Nuovo, Usato, Km&Oslash;</h2>
          </div>
          <div id="page-header-contacts" class="large-3 small-12 columns">
            <a href="/contacts" class="alert radius button">Dove Siamo</a>
          </div>
        </header>

        <section id="page-content" class="row">
          <div class="small-12 columns" ng-view ng-animate="'view-animation'">
          	<?php
          		if ( $this->agent->is_robot() ) {
          			$segment = $this->uri->segment(1);
          			if ( !$segment ) $segment = 'carlist';
          			$this->load->view($segment, array( 'cars' => $cars ));
          		}
          	?>
          </div>
        </section>
      </div>

      <footer id="site-footer" class="small-12 columns">
        Copyright &copy; 2013 Union Car S.r.l
      </footer>
    </div>

    <script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.2/jquery.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/select2/3.4.0/select2.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/flexslider/2.1/jquery.flexslider-min.js"></script>
    <script src="<?php echo base_url(APPPATH.'/scripts/scripts.js'); ?>"></script>
    <script type="text/javascript">
    angular.module('UnionCarWebsiteApp').constant('conf', {
    	baseUrl: "<?php echo rtrim(base_url(), '/'); ?>",
    	appUrl: "<?php echo base_url(APPPATH); ?>"
    });
    </script>

    <!-- Google Analytics: change UA-XXXXX-X to be your site's ID. -->
    <script>
      var _gaq=[['_setAccount','UA-XXXXX-X'],['_trackPageview']];
      (function(d,t){var g=d.createElement(t),s=d.getElementsByTagName(t)[0];
      g.src=('https:'==location.protocol?'//ssl':'//www')+'.google-analytics.com/ga.js';
      s.parentNode.insertBefore(g,s)}(document,'script'));
    </script>
  </body>
</html>
