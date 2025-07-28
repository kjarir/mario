import { animate, svg, stagger } from 'animejs';
import { useEffect } from 'react';

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
    </div>
  );
}

export default Index;
