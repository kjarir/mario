import { animate, svg, stagger } from 'animejs';
import { useEffect } from 'react';
import ThreeDHoverGallery from '../components/lightswind/3d-hover-gallery';
import ScrollReveal from '../components/lightswind/scroll-reveal';

function Index() {
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
      <div
        className="w-full h-screen flex items-center justify-center"
        style={{ padding:"15%" }}
      >
        <svg
          viewBox="0 0 304 112"
          xmlns="http://www.w3.org/2000/svg"
          className="w-100 h-100 text-white"
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
      </div>
      
      {/* 3D Hover Gallery Section */}
      <div className="bg-white py-16">
        <div className="container mx-auto px-4">
          
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

      {/* Scroll Reveal Section */}
      <div className="bg-gray-100 h-screen flex items-center justify-center">
        <div className="container mx-auto px-4">
          <div className="w-full space-y-8">
            <ScrollReveal
              size="2xl"
              align="center"
              variant="primary"
              enableBlur={true}
              staggerDelay={0.1}
              duration={1.2}
              textClassName="text-7xl md:text-9xl lg:text-[10rem] xl:text-[15rem] font-bold leading-tight text-center"
            >
              Introducing to you Dr.Mario
            </ScrollReveal>

            <ScrollReveal
              size="xl"
              align="center"
              variant="muted"
              enableBlur={true}
              staggerDelay={0.05}
              duration={0.8}
              textClassName="text-4xl md:text-6xl lg:text-7xl xl:text-8xl font-medium leading-relaxed text-center"
            >
              Revolutionary AI-powered diabetic retinopathy detection through retinal imaging
            </ScrollReveal>

            <ScrollReveal
              size="2xl"
              align="center"
              variant="default"
              enableBlur={true}
              staggerDelay={0.08}
              duration={1.0}
              textClassName="text-6xl md:text-7xl lg:text-8xl xl:text-9xl font-bold leading-tight text-center"
            >
              Early detection saves vision and lives
            </ScrollReveal>

            <ScrollReveal
              size="xl"
              align="center"
              variant="accent"
              enableBlur={true}
              staggerDelay={0.06}
              duration={0.9}
              textClassName="text-5xl md:text-6xl lg:text-7xl xl:text-8xl font-semibold leading-relaxed text-center"
            >
              Advanced machine learning algorithms analyze retinal scans with 99% accuracy
            </ScrollReveal>

            <ScrollReveal
              size="2xl"
              align="center"
              variant="primary"
              enableBlur={true}
              staggerDelay={0.12}
              duration={1.1}
              textClassName="text-6xl md:text-7xl lg:text-8xl xl:text-9xl font-bold leading-tight text-center"
            >
              Transforming healthcare through intelligent retinal analysis
            </ScrollReveal>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Index;
