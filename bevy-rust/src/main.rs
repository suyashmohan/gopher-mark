use bevy::{
    diagnostic::{DiagnosticsStore, FrameTimeDiagnosticsPlugin},
    prelude::*,
};
use rand::Rng;

const SCREEN_WIDTH: u32 = 1920;
const SCREEN_HEIGHT: u32 = 1080;
const GOPHER_COUNT: usize = 1000;

fn main() {
    App::new()
        .add_plugins((
            DefaultPlugins.set(WindowPlugin {
                primary_window: Some(Window {
                    title: "gopher mark".into(),
                    resolution: (SCREEN_WIDTH, SCREEN_HEIGHT).into(),
                    present_mode: bevy::window::PresentMode::Fifo,
                    ..default()
                }),
                ..default()
            }),
            FrameTimeDiagnosticsPlugin::default(),
        ))
        .add_systems(Startup, setup)
        .add_systems(Update, (gopher_movement, spwan_gophers, update_fps_text))
        .run();
}

#[derive(Component)]
struct Gopher {
    pub velocity: Vec2,
}

#[derive(Resource)]
struct SpawnTimer {
    counter: usize,
}

#[derive(Resource)]
struct GopherAssets {
    handle: Handle<Image>,
}

fn setup(mut commands: Commands, asset_server: Res<AssetServer>) {
    let gopher_image = asset_server.load("gopher.png");
    let mut rng = rand::rng();

    commands.insert_resource(GopherAssets {
        handle: gopher_image.clone(),
    });

    commands.spawn(Camera2d);

    commands.spawn((
        Text::new("FPS: "),
        Node {
            position_type: PositionType::Absolute,
            top: Val::Px(10.0),
            left: Val::Px(10.0),
            ..default()
        },
    ));

    for _ in 0..GOPHER_COUNT {
        let velocity = Vec2 {
            x: rng.random_range(-100.0..100.0),
            y: rng.random_range(-100.0..100.0),
        };
        commands.spawn((
            Gopher { velocity },
            Sprite::from_image(gopher_image.clone()),
            Transform::from_xyz(0., 0., 0.),
        ));
    }

    commands.insert_resource(SpawnTimer {
        counter: GOPHER_COUNT,
    });
}

fn spwan_gophers(
    mut commands: Commands,
    mut spwan_timer: ResMut<SpawnTimer>,
    assets: Res<GopherAssets>,
    inputs: Res<ButtonInput<KeyCode>>,
) {
    if inputs.just_released(KeyCode::Space) {
        let mut rng = rand::rng();

        let gophers_to_spawn: Vec<_> = (0..GOPHER_COUNT)
            .map(|_| {
                let velocity = Vec2 {
                    x: rng.random_range(-100.0..100.0),
                    y: rng.random_range(-100.0..100.0),
                };
                (
                    Gopher { velocity },
                    Sprite::from_image(assets.handle.clone()),
                    Transform::from_xyz(0., 0., 0.),
                )
            })
            .collect();

        commands.spawn_batch(gophers_to_spawn);
        spwan_timer.counter += GOPHER_COUNT;
    }
}

fn gopher_movement(time: Res<Time>, mut query: Query<(&mut Gopher, &mut Transform)>) {
    let half_width = SCREEN_WIDTH as f32 / 2.0;
    let half_height = SCREEN_HEIGHT as f32 / 2.0;

    let delta = time.delta_secs();

    query
        .par_iter_mut()
        .for_each(|(mut gopher, mut transform)| {
            if transform.translation.x < -half_width || transform.translation.x > half_width {
                gopher.velocity.x *= -1.0;
            }
            if transform.translation.y < -half_height || transform.translation.y > half_height {
                gopher.velocity.y *= -1.0;
            }
            transform.translation.x += gopher.velocity.x * delta;
            transform.translation.y += gopher.velocity.y * delta;
        });
}

fn update_fps_text(
    diagnositcs: Res<DiagnosticsStore>,
    spawn_timer: Res<SpawnTimer>,
    mut query: Query<&mut Text>,
) {
    if let Some(fps) = diagnositcs.get(&FrameTimeDiagnosticsPlugin::FPS) {
        if let Some(value) = fps.smoothed() {
            for mut text in &mut query {
                text.0 = format!("FPS: {:.1} - Gophers: {}", value, spawn_timer.counter);
            }
        }
    }
}
