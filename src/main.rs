#[macro_use]
extern crate glium;

use glium::Surface;

fn main() {

    // Create the event loop
    let event_loop = winit::event_loop::EventLoopBuilder::new().build().expect("event loop building");


    // Create the window
    let (window, display) = glium::backend::glutin::SimpleWindowBuilder::new().build(&event_loop);

    // deltaT increase over time
    let mut delta_t: f32 = 0.0;

    let _ = event_loop.run(move |event, window_target| {
    match event {
        winit::event::Event::WindowEvent { event, .. } => match event {
            // Close window
            winit::event::WindowEvent::CloseRequested => window_target.exit(),

            // Resize window event
            winit::event::WindowEvent::Resized(window_size) => {
                display.resize(window_size.into());
            },

            // Draw loop
            winit::event::WindowEvent::RedrawRequested => {

                // Increase delta time over time
                delta_t += 0.02;

                // Sine offset
                let x_off = delta_t.sin() * 0.5;
                
                // Draw target
                let mut target = display.draw();

                // Set background colour
                target.clear_color(0.02, 0.02, 0.02, 1.0);


                #[derive(Copy, Clone)]
                struct Vertex {
                    position: [f32; 2],
                }

                // Triangle vertices
                let vertex1 = Vertex { position: [-0.5 + x_off, -0.5] };
                let vertex2 = Vertex { position: [ 0.0 + x_off,  0.5] };
                let vertex3 = Vertex { position: [ 0.5 + x_off, -0.25] };

                // Create the shape array from the vertices
                let shape = vec![vertex1, vertex2, vertex3];
                implement_vertex!(Vertex, position);

                // Create a vertex buffer using the shape
                let vertex_buffer = glium::VertexBuffer::new(&display, &shape).unwrap();

                // Indices
                let indices = glium::index::NoIndices(glium::index::PrimitiveType::TrianglesList);

                // Vertex shader
                // GLSL language
                // --- A vertex shader is a small program that tells the GPU where to draw vertices
                // in relation to the screen coordinates
                let vertex_shader_src = r#"
                    #version 140

                    in vec2 position;

                    void main() {
                        gl_Position = vec4(position, 0.0, 1.0);
                    }
                "#;

                // Fragment shader
                // A fragment shader is a small program whose purpose is to tell the GPU what color each pixel should be
                let fragment_shader_src = r#"
                    #version 140

                    out vec4 color;

                    void main() {
                        color = vec4(1.0, 0.0, 0.0, 1.0);
                    }
                "#;


                // Send the shaders to Glium
                let program = glium::Program::from_source(&display, vertex_shader_src, fragment_shader_src, None).unwrap();

                // Draw the triangle
                target.draw(&vertex_buffer, &indices, &program, &glium::uniforms::EmptyUniforms,
                        &Default::default()).unwrap();


                // Send to the window
                target.finish().unwrap();

            },
            _ => (),
        },
        winit::event::Event::AboutToWait => {
            // Good practice to redraw
            window.request_redraw();
        },
        _ => (),
    };
});

}
