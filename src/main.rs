mod manager;
use gtk::cairo;
use gtk::gdk;
use gtk::glib;
use gtk::prelude::*;
use gtk4 as gtk;
use std::cell::RefCell;
use std::time::Duration;
use std::time::SystemTime;
use std::{cell::Cell, rc::Rc};

use crate::manager::Manager;
fn main() {
    let application = gtk::Application::builder()
        .application_id("com.github.gtk-rs.examples.dialog_async")
        .build();

    application.connect_activate(build_ui);

    application.run();
}

fn build_ui(application: &gtk::Application) {
    let window = gtk::ApplicationWindow::builder()
        .application(application)
        .title("Dialog Async")
        .default_width(350)
        .default_height(70)
        .visible(true)
        .build();
    let area = gtk::DrawingArea::new();
    const FONT_SIZE: f64 = 20.0;
    let manager = Rc::new(RefCell::new(Manager::new()));
    // let mut text = Rc::new(RefCell::new(String::new()));
    let line_height = FONT_SIZE + 10.0;
    let eventctl = gtk::EventControllerKey::new();

    eventctl.connect_key_pressed(
        glib::clone!(@strong manager, @strong area => move |_eventctl, keyval, _keycode, _state| {
                let mut t = manager.borrow_mut();
                t.read_key(keyval);
            area.queue_draw();
            glib::signal::Propagation::Proceed
        }),
    );
    eventctl.connect_key_released(glib::clone!(@strong manager => move |_eventctl, _keyval, _keycode, _state| {
        // Check if one of alt, shift, ctrl, super is clicked and set it on the manager
        let mut t = manager.borrow_mut();
        t.alt = false;
        t.ctrl = false;
        t.shift = false;
        t.superkey = false;

    }));
    area.set_draw_func( glib::clone!(@strong line_height,  => move|area, cr, _w, _h| {
        cr.set_source_rgb(0.15625, 0.1640625, 0.2109375);
        cr.paint().expect("Cannot paint");
        cr.set_antialias(cairo::Antialias::Subpixel);
        cr.set_font_size(FONT_SIZE);
        // Highlight the current line
        cr.set_source_rgb(0.2, 0.2, 0.2);
        let screen_width = area.width();
        cr.rectangle(0.0, (manager.borrow().y as f64 + 0.2) * line_height, screen_width as f64, line_height + 0.5);
        cr.fill().expect("Cannot fill");
        

        // Show curser
        cr.set_source_rgb(0.4, 0.4, 0.4);
        cr.rectangle(50. + (manager.borrow().x as f64 * 11.), (manager.borrow().y as f64 + 1.) * line_height, 10., 2.0);
        cr.fill().expect("Cannot fill");
        let man = manager.borrow();
        let lines = man.text.split("\n").collect::<Vec<&str>>();
        if lines.len() != 0 {
            for (i, line) in lines.into_iter().enumerate() {
                cr.select_font_face("Mono", cairo::FontSlant::Normal, cairo::FontWeight::Normal);
                cr.move_to(15.0, (i as f64 + 1.0) * line_height);
                cr.set_source_rgb(0.256625, 0.27734375, 0.3515625);
                cr.show_text((i + 1).to_string().as_str()).expect("Cant show text!");
                let mut j = 0;
                for char in line.to_string().chars() {
                    cr.select_font_face("Mono", cairo::FontSlant::Normal, cairo::FontWeight::Normal);
                    cr.move_to(50.+(j as f64 * 11.), (i as f64 + 1.) * line_height);
                    cr.set_source_rgb(0.96875, 0.96875, 0.9453125);
                    cr.show_text(&char.to_string()).expect("Cannot show text");
                    j+=1;
                }
            }
        }
        area.queue_draw();
    }));
    window.add_controller(eventctl);
    window.set_child(Some(&area));
    window.present()
}
