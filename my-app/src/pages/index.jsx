import { animate, svg, stagger } from 'animejs';
import { useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import ThreeDHoverGallery from '../components/lightswind/3d-hover-gallery';
import VariableProximity from '../ReactBits/VariableProximity/VariableProximity';

function Index() {
  const containerRef = useRef(null);

  useEffect(() => {
    animate(svg.createDrawable('.line'), {
      draw: ['0 0', '0 1', '1 1'],
      ease: 'inOutQuad',
      duration: 2000,
      delay: stagger(100),
      loop: true
    });
  }, []);

  return (
    <div className="min-h-screen" style={{ backgroundColor: '#D9272C' }}>
      {/* Navigation Header */}
      <div className="absolute top-0 left-0 right-0 z-50 p-6">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="text-white text-2xl font-bold">Dr. Mario</div>
          <div className="flex space-x-4">
            <Link
              to="/login"
              className="text-white hover:text-gray-200 px-4 py-2 rounded-md transition-colors"
            >
              Sign In
            </Link>
            <Link
              to="/register"
              className="bg-white text-red-600 hover:bg-gray-100 px-4 py-2 rounded-md font-medium transition-colors"
            >
              Get Started
            </Link>
          </div>
        </div>
      </div>

      {/* Hero Section */}
      <div
        className="w-full h-screen flex items-center justify-center"
        style={{ padding:"15%" }}
      >
        <div className="text-center">
          <svg
            viewBox="0 0 304 112"
            xmlns="http://www.w3.org/2000/svg"
            className="w-100 h-100 text-white mx-auto mb-8"
          >
            <g
              stroke="currentColor"
              fill="none"
              fillRule="evenodd"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
            >
              <path className="line" d="M8 90V22h17l17 34 17-34h17v68h-17V56l-17 34-17-34v34H8z" />
              <path className="line" d="M125 90V56.136C124.66 46.48 117.225 39 108 39c-9.389 0-17 7.611-17 17s7.611 17 17 17h8.5v17H108c-18.778 0-34-15.222-34-34s15.222-34 34-34c18.61 0 33.433 14.994 34 33.875V90h-17z" />
              <path className="line" d="M160 90V22h17c9.389 0 17 7.611 17 17s-7.611 17-17 17h-8.5v34H160z" />
              <path className="line" d="M179 56l21 34h-18l-13-21V56h10z" />
              <path className="line" d="M205 90V22h17v68h-17z" />
              <path className="line" d="M274 90c-18.778 0-34-15.222-34-34s15.222-34 34-34c18.778 0 34 15.222 34 34s-15.222 34-34 34z" />
              <path className="line" d="M274 73a17 17 0 1 1 0-34 17 17 0 0 1 0 34z" />
            </g>
          </svg>
          
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-6">
            AI-Powered Retinal Imaging
          </h1>
          <p className="text-xl text-white/90 mb-8 max-w-2xl mx-auto">
            Advanced diabetic retinopathy detection using deep learning and medical image analysis
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/register"
              className="bg-white text-red-600 hover:bg-gray-100 px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
            >
              Start Free Trial
            </Link>
            <Link
              to="/login"
              className="border-2 border-white text-white hover:bg-white hover:text-red-600 px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
            >
              Sign In
            </Link>
          </div>
        </div>
      </div>
      
      {/* 3D Hover Gallery Section */}
      <div className="bg-white py-16">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12 text-gray-900">
            Advanced Image Analysis
          </h2>
          <ThreeDHoverGallery 
            images={[
              "https://images.pexels.com/photos/26797335/pexels-photo-26797335/free-photo-of-scenic-view-of-mountains.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/12194487/pexels-photo-12194487.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32423809/pexels-photo-32423809/free-photo-of-aerial-view-of-kayaking-at-robberg-south-africa.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32296519/pexels-photo-32296519/free-photo-of-rocky-coastline-of-cape-point-with-turquoise-waters.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32396739/pexels-photo-32396739/free-photo-of-serene-motorcycle-ride-through-bamboo-grove.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32304900/pexels-photo-32304900/free-photo-of-scenic-view-of-cape-town-s-twelve-apostles.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32437034/pexels-photo-32437034/free-photo-of-fisherman-holding-freshly-caught-red-drum-fish.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
              "https://images.pexels.com/photos/32469847/pexels-photo-32469847/free-photo-of-deer-drinking-from-natural-water-source-in-wilderness.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
            ]}
            itemWidth={12}
            itemHeight={20}
            gap={1.2}
            perspective={50}
            hoverScale={15}
            transitionDuration={1.25}
            backgroundColor="#ffffff"
            grayscaleStrength={1}
            brightnessLevel={0.5}
            activeWidth={45}
            rotationAngle={35}
            zDepth={10}
            enableKeyboardNavigation={true}
            autoPlay={false}
            autoPlayDelay={3000}
            onImageClick={(index, image) => console.log(`Clicked image ${index}: ${image}`)}
            onImageHover={(index, image) => console.log(`Hovered image ${index}: ${image}`)}
            onImageFocus={(index, image) => console.log(`Focused image ${index}: ${image}`)}
          />
        </div>
      </div>

      {/* Variable Proximity Section */}
      <div 
        ref={containerRef}
        className="bg-gray-100 h-screen flex items-center justify-center p-10 cursor-pointer"
      >
        <div className="w-full max-w-6xl">
          <div className="text-[3.2rem] md:text-[3rem] lg:text-[4rem] xl:text-[4rem] font-light leading-tight text-left">
            <VariableProximity
              label="Dr.Mario is an AI-powered retinal imaging system designed to detect early signs of Diabetic Retinopathy (DR). Leveraging deep learning and medical image analysis, it aims to assist ophthalmologists by providing fast, accurate, and non-invasive diagnosis from fundus images enabling timely intervention and reducing the risk of vision loss in diabetic patients."
              fromFontVariationSettings="'wght' 100"
              toFontVariationSettings="'wght' 900"
              containerRef={containerRef}
              radius={150}
              falloff="linear"
              className="text-black cursor-pointer"
            />
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-xl font-bold mb-4">Dr. Mario</h3>
              <p className="text-gray-400">
                AI-powered retinal imaging for early diabetic retinopathy detection.
              </p>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Features</h4>
              <ul className="space-y-2 text-gray-400">
                <li>AI Detection</li>
                <li>Image Analysis</li>
                <li>Secure Storage</li>
                <li>Real-time Results</li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">For</h4>
              <ul className="space-y-2 text-gray-400">
                <li>Patients</li>
                <li>Doctors</li>
                <li>Clinics</li>
                <li>Hospitals</li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Get Started</h4>
              <div className="space-y-2">
                <Link
                  to="/register"
                  className="block text-gray-400 hover:text-white transition-colors"
                >
                  Create Account
                </Link>
                <Link
                  to="/login"
                  className="block text-gray-400 hover:text-white transition-colors"
                >
                  Sign In
                </Link>
              </div>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2024 Dr. Mario. All rights reserved.</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Index;
