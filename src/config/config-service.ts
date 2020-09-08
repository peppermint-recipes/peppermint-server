import { Injectable } from "@nestjs/common";
import { TypeOrmOptionsFactory, TypeOrmModuleOptions } from "@nestjs/typeorm";

import * as dotenv from "dotenv";
import * as fs from 'fs';
import * as path from 'path';

import { Recipe } from "../recipe/recipe-model";

// Based on https://stackoverflow.com/questions/54361685/nestjs-typeorm-configuration-using-env-files/54364907#54364907
@Injectable()
export class ConfigService implements TypeOrmOptionsFactory {
    private readonly envConfig: dotenv.DotenvParseOutput;
    private readonly environment: string;

    constructor() {
        this.environment = process.env.NODE_ENV || "local"
        const fileName = path.join(process.cwd(), `${this.environment}.env`);
        this.envConfig = dotenv.parse(
            fs.readFileSync(fileName)
        );
    }

    createTypeOrmOptions(): TypeOrmModuleOptions {
        return {
            type: "postgres",
            host: this.envConfig.DBHOST,
            port: Number(this.envConfig.DBPORT),
            username: this.envConfig.DBUSER,
            password: this.envConfig.DBPASSWORD,
            database: this.envConfig.DBDATABASE,
            entities: [Recipe],
            synchronize: true,
        };
    }
}
