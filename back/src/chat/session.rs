use std::time::{Duration, Instant};

use actix::{Actor, ActorContext, AsyncContext, StreamHandler};
use actix_web_actors::ws;

/// How often heartbeat pings are sent
const HEARTBEAT_INTERVAL: Duration = Duration::from_secs(5);

/// How long before lack of client response causes a timeout
const CLIENT_TIMEOUT: Duration = Duration::from_secs(10);

/// Define Websocket actor
pub struct WsChatSession {
    pub hb: Instant,

    #[allow(dead_code)] // TODO: will work soon
    room_id: u32,
}

impl Actor for WsChatSession {
    type Context = ws::WebsocketContext<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        self.hb(ctx);
    }
}

pub fn new(room_id: u32) -> WsChatSession {
    WsChatSession {
        hb: Instant::now(),
        room_id,
    }
}


impl WsChatSession {
    fn hb(&self, ctx: &mut ws::WebsocketContext<Self>) {
        ctx.run_interval(HEARTBEAT_INTERVAL, |act, ctx| {
            // check client heartbeats
            if Instant::now().duration_since(act.hb) > CLIENT_TIMEOUT {
                // heartbeat timed out
                println!("Websocket Client heartbeat failed, disconnecting!");

                ctx.stop();
                return;
            }

            ctx.text("ping");
        });
    }
}

/// Handler for ws::Message message
impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for WsChatSession {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(ws::Message::Text(text)) => {
                if let Err(e) = String::from_utf8(text.as_bytes().to_vec()) {
                    let reason = ws::CloseReason{ code: ws::CloseCode::Error, description: Some(format!("can't parse message: {e}")) };
                    ctx.close(Some(reason));
                }

                let str = String::from_utf8(text.as_bytes().to_vec()).unwrap();
                match str.as_str() {
                    "ping" => {
                        self.hb = Instant::now();
                    }
                    "pong" => {
                        self.hb = Instant::now();
                    }
                    _ => { ctx.text(text) }
                }
            }

            Ok(ws::Message::Close(_)) => {}
            Ok(ws::Message::Nop) => {}

            Err(_) => {}
            _ => (),
        }
    }
}

