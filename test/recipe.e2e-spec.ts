import { getConnection } from "typeorm";
import { INestApplication } from "@nestjs/common";
import { Test } from "@nestjs/testing";
import { TypeOrmModule } from "@nestjs/typeorm";

import * as request from "supertest";

import { ConfigModule } from "../src/config/config-module";
import { ConfigService } from "../src/config/config-service";
import { Recipe } from "../src/recipe/recipe-model";
import { RecipeModule } from "../src/recipe/recipe-module";

describe("RecipeController (/api/recipes) (e2e)", () => {
    let app: INestApplication;

    // Based on https://github.com/nestjs/nest/issues/409
    const cleanupDatabase = async () => {
        await getConnection().synchronize(true);
    };

    beforeAll(async () => {
        const moduleFixture = await Test.createTestingModule({
            imports: [
                RecipeModule,
                TypeOrmModule.forRootAsync({
                    imports: [ConfigModule],
                    useExisting: ConfigService,
                }),
            ],
        }).compile();

        app = moduleFixture.createNestApplication();
        await app.init();
    });

    beforeEach(async () => {
        await cleanupDatabase();
    });

    it("should create a valid recipe)", async () => {
        expect.assertions(2);

        const recipe = new Recipe({
            name: "Test-Recipe",
            servings: "2",
            activeTime: "5",
            totalTime: "10",
            ingredients: "Banana",
            instructions: "Do stuff",
        });
        const createRecipeResponse = await request(app.getHttpServer())
            .post("/api/recipes")
            .send(recipe);
        const createdTag = createRecipeResponse.body;
        expect(createRecipeResponse.status).toStrictEqual(201);
        expect(createdTag.name).toStrictEqual("Test-Recipe");
    });

    it("should not create an invalid recipe)", async () => {
        expect.assertions(1);

        const recipe: Partial<Recipe> = {
            name: "Test-123",
        };
        const createTagRecipe = await request(app.getHttpServer())
            .post("/api/recipes")
            .send(recipe);
        expect(createTagRecipe.status).toStrictEqual(500);
    });

    it("should GET all recipies (empty array)", async () => {
        expect.assertions(2);

        const response = await request(app.getHttpServer()).get("/api/recipes");

        expect(response.status).toStrictEqual(200);
        expect(response.body).toStrictEqual({ recipes: [] });
    });

    it("should DELETE a recipe by id)", async () => {
        expect.assertions(5);

        const recipe = new Recipe({
            name: "Test-Recipe",
            servings: "2",
            activeTime: "5",
            totalTime: "10",
            ingredients: "Banana",
            instructions: "Do stuff",
        });

        const createRecipeResponse = await request(app.getHttpServer())
            .post("/api/recipes")
            .send(recipe);

        const createdRecipe = createRecipeResponse.body;

        expect(createRecipeResponse.status).toStrictEqual(201);
        expect(createdRecipe.name).toStrictEqual("Test-Recipe");

        const deleteRecipeResponse = await request(app.getHttpServer()).delete(
            "/api/recipes/" + createdRecipe.id
        );

        expect(deleteRecipeResponse.status).toStrictEqual(200);

        const response = await request(app.getHttpServer()).get("/api/recipes");
        expect(response.status).toStrictEqual(200);
        expect(response.body).toStrictEqual({ recipes: [] });
    });

    it("should get a recipe by id)", async () => {
        expect.assertions(4);

        const recipe = new Recipe({
            name: "Test-Recipe",
            servings: "2",
            activeTime: "5",
            totalTime: "10",
            ingredients: "Banana",
            instructions: "Do stuff",
        });

        const createRecipeResponse = await request(app.getHttpServer())
            .post("/api/recipes")
            .send(recipe);

        const createdRecipe = createRecipeResponse.body;
        expect(createRecipeResponse.status).toStrictEqual(201);
        expect(createdRecipe.name).toStrictEqual("Test-Recipe");

        const response = await request(app.getHttpServer()).get(
            "/api/recipes/" + createdRecipe.id
        );

        expect(response.status).toStrictEqual(200);
        expect(response.body).toStrictEqual(createdRecipe);
    });

    it("should update a tag by id)", async () => {
        expect.assertions(5);

        const recipe = new Recipe({
            name: "Test-Recipe",
            servings: "2",
            activeTime: "5",
            totalTime: "10",
            ingredients: "Banana",
            instructions: "Do stuff",
        });

        const createTagResponse = await request(app.getHttpServer())
            .post("/api/recipes")
            .send(recipe);

        const createdRecipe = createTagResponse.body;
        expect(createTagResponse.status).toStrictEqual(201);
        expect(createdRecipe.name).toStrictEqual("Test-Recipe");

        const updatedRecipe = new Recipe({
            name: "Test-1337",
            servings: "5",
            activeTime: "15",
            totalTime: "100",
            ingredients: "Banana, Oil",
            instructions: "Do stuff",
        });

        const response = await request(app.getHttpServer())
            .put("/api/recipes/" + createdRecipe.id)
            .send(updatedRecipe);
        expect(response.status).toStrictEqual(200);
        expect(response.body.name).toStrictEqual("Test-1337");
        expect(response.body.ingredients).toStrictEqual("Banana, Oil");
    });
});
