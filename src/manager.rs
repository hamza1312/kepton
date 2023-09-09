#[derive(Clone)]
pub struct Manager {
    pub text: String,
    pub x: usize,
    pub y: usize,
    pub alt: bool,
    pub ctrl: bool,
    pub shift: bool,
    pub superkey: bool,
}

impl Manager {
    pub fn new() -> Self {
        Self {
            text: String::new(),
            x: 0,
            y: 0,
            alt: false,
            ctrl: false,
            shift: false,
            superkey: false,
        }
    }
    fn get_offset(&self) -> usize {
        let lines: Vec<&str> = self.text.split('\n').collect();
        let mut offset = 0;
        for (i, line) in lines.iter().enumerate() {
            if self.y > i {
                offset += line.len() + 1;
            }
            if self.y == i {
                offset += self.x;
                break;
            }
        }
        offset
    }
    pub fn bksp_char(&mut self) {
        let offset = self.get_offset();
        if offset > 0 {
            self.text = format!("{}{}", &self.text[..offset - 1], &self.text[offset..]);
            if self.x > 0 {
                self.x -= 1;
            }
        }
    }
    pub fn read_key(&mut self, key: gtk4::gdk::Key) {
        if key.name().is_some() {
            match key.name().unwrap().to_string().as_str() {
                "Return" => {
                    let offset = self.get_offset();
                    self.text = format!("{}\n{}", &self.text[..offset], &self.text[offset..]);
                    self.y += 1;
                    self.x = 0;
                }
                // Tab
                "Tab" => {
                    let offset = self.get_offset();
                    self.text = format!("{}    {}", &self.text[..offset], &self.text[offset..]);
                    self.x += 4;
                }
                // Arrow keys:
                "Up" => {
                    if self.y > 0 {
                        self.y -= 1;
                        // Get the line and use the length of the line to set the x position
                        let lines: Vec<&str> = self.text.split('\n').collect();
                        let line = lines[self.y].to_string();
                        self.x = line.len();
                    }
                }
                "Down" => {
                    if self.text.split('\n').collect::<Vec<_>>().len() - 1 > self.y {
                        self.y += 1;
                        // Get the line and use the length of the line to set the x position
                        let lines: Vec<&str> = self.text.split('\n').collect();
                        if self.y > lines.len() - 1 {
                            self.y = lines.len() - 1;
                        }
                        let line = lines[self.y].to_string();
                        self.x = line.len();
                    }
                }
                "Right" => {
                    if self.ctrl {
                        let lines: Vec<&str> = self.text.split('\n').collect();
                        let line = lines[self.y].to_string();
                        let mut chars = line.chars();
                        let mut i = 0;
                        while i < self.x {
                            chars.next();
                            i += 1;
                        }
                        let mut word = String::new();
                        let mut j = 0;
                        while let Some(c) = chars.next() {
                            if c == ' ' {
                               
                                break;
                            }
                            word.push(c);
                            j += 1;
                        }
                        self.x += j;
                    } else {
                        // Check if the line length is greater than the x position
                        let lines: Vec<&str> = self.text.split('\n').collect();
                        if lines.get(self.y).is_some() {
                            if lines[self.y].len() > self.x {
                                self.x += 1;
                            }
                        }
                    }
                }
                "Left" => {
                    if self.ctrl {
                        let lines: Vec<&str> = self.text.split('\n').collect();
                        let line = lines[self.y].to_string();
                        let mut chars = line.chars();
                        let mut i = 0;
                        while i < self.x {
                            chars.next();
                            i += 1;
                        }
                        let mut word = String::new();
                        let mut j = 0;
                        while let Some(c) = chars.next() {
                            if c == ' ' {
                                break;
                            }
                            word.push(c);
                            j += 1;
                        }
                        self.x -= j;
                    } else {
                        if self.x > 0 {
                            self.x -= 1;
                        }
                    }
                }
                "BackSpace" => {
                    let text_clone = self.text.clone();
                    let lines = text_clone.split('\n').collect::<Vec<_>>();
                    if lines.get(self.y).is_some() {
                        let line = lines[self.y];
                        // If there is nothing on the line, just delete it
                        if line.len() == 0 {
                            if self.y != 0 {
                                self.text = format!(
                                    "{}{}",
                                    &self.text[..self.get_offset() - 1],
                                    &self.text[self.get_offset()..]
                                );
                                self.y -= 1;
                                self.x = lines[self.y].len();
                            }
                        } else {
                            self.bksp_char();
                        }
                    }
                }
                // alt, ctrl, shift, super
                "Alt_L" => {
                    self.alt = true;
                }
                "Alt_R" => {
                    self.alt = true;
                }
                "Control_L" => {
                    self.ctrl = true;
                }
                "Control_R" => {
                    self.ctrl = true;
                }
                "Shift_L" => {
                    self.shift = true;
                }
                "Shift_R" => {
                    self.shift = true;
                }
                "Super_L" => {
                    self.superkey = true;
                }
                "Super_R" => {
                    self.superkey = true;
                }

                _ => {
                    let key_unicode = key.to_unicode().unwrap_or('\0');
                    if key_unicode as u32 >= 32 && key_unicode as u32 <= 127 {
                        // push to the text depending on the x and y position of the cursor
                        let mut autocomplete = String::new();
                        if key_unicode == '(' {
                            autocomplete = ")".into()
                        }
                        if key_unicode == '[' {
                            autocomplete = "]".into()
                        }
                        if key_unicode == '{' {
                            autocomplete = "}".into()
                        }
                        if key_unicode == '"' {
                            autocomplete = "\"".into()
                        }
                        if key_unicode == '\'' {
                            autocomplete = "'".into()
                        }
                        let offset = self.get_offset();
                        self.text = format!(
                            "{}{}{}{}",
                            &self.text[..offset],
                            key_unicode,
                            autocomplete,
                            &self.text[offset..]
                        );

                        self.x += 1
                    }
                }
            }
        }
    }
}
